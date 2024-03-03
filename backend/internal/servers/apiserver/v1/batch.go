package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/common"
	"github.com/darylhjd/oams/backend/pkg/goroutines"
	"github.com/darylhjd/oams/backend/pkg/to"
)

const (
	maxBatchPostParseMemory     = 32 << 20
	maxBatchPostGoRoutines      = 10
	multipartFormBatchFileIdent = "batch-attachments"
	multipartFormBatchWeekIdent = "start-week"
)

func (v *APIServerV1) batch(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodPost:
		resp = v.batchPost(r)
	case http.MethodPut:
		resp = v.batchPut(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type batchPostResponse struct {
	response
	batchPutRequest
}

// batchPost processes a file and returns the corresponding PUT request that can be created
// from it. It does not process (create, delete, etc...) any of the entities.
func (v *APIServerV1) batchPost(r *http.Request) apiResponse {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart") {
		return newErrorResponse(http.StatusUnsupportedMediaType, "a multipart request body is required")
	}

	resp, err := v.processBatchPostRequest(r)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process batch file(s)")
	}

	return resp
}

func (v *APIServerV1) processBatchPostRequest(r *http.Request) (apiResponse, error) {
	err := r.ParseMultipartForm(maxBatchPostParseMemory)
	if err != nil {
		return batchPostResponse{}, err
	}

	var startWeek int64
	startWeekField := r.MultipartForm.Value[multipartFormBatchWeekIdent]
	if len(startWeekField) != 1 {
		return newErrorResponse(http.StatusBadRequest, "unexpected start week field format"), nil
	} else if startWeek, err = to.Int64(startWeekField[0]); err != nil {
		return newErrorResponse(http.StatusBadRequest, "start week field is not an integer"), nil
	}

	limiter := goroutines.NewLimiter(maxBatchPostGoRoutines)
	saveRes := sync.Map{}
	for _, header := range r.MultipartForm.File[multipartFormBatchFileIdent] {
		header := header // Required for go routine to point to different file for each loop.
		limiter.Do(func() {
			var data common.BatchData

			file, err := header.Open()
			if err != nil {
				// Save value as error type. This is an internal error.
				saveRes.Store(&data, err)
				return
			}
			defer func() {
				_ = file.Close()
			}()

			data, err = common.ParseBatchFile(header.Filename, int(startWeek), file)
			if err != nil {
				// Save as string type. This is a request error.
				saveRes.Store(&data, err.Error())
				return
			}

			saveRes.Store(&data, nil)
		})
	}

	limiter.Wait()

	var (
		errResp       errorResponse
		okResp        batchPostResponse
		isErrResponse bool
	)
	saveRes.Range(func(key, value any) bool {
		data, ok := key.(*common.BatchData)
		if !ok {
			err = errors.New("type assertion failed when processing batch data")
			return false
		}

		switch t := value.(type) {
		case error:
			err = t
			return false
		case string:
			isErrResponse = true
			errResp = newErrorResponse(http.StatusBadRequest, t)
			return false
		case nil:
			okResp.Batches = append(okResp.Batches, *data)
			return true
		default:
			err = errors.New("type assertion failed when processing batch result")
			return false
		}
	})

	switch {
	case err != nil:
		return batchPostResponse{}, err
	case isErrResponse:
		return errResp, nil
	default:
		okResp.response = response{true, http.StatusAccepted}
		return okResp, nil
	}
}

type batchPutRequest struct {
	Batches []common.BatchData `json:"batches"`
}

type batchPutResponse struct {
	response
	ClassIDs []int64 `json:"class_ids"`
}

// batchPut is the handler that does the actual processing of the entities.
func (v *APIServerV1) batchPut(r *http.Request) apiResponse {
	var req batchPutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	resp, err := v.processBatchPutRequest(r, req)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process batch put database action")
	}

	return resp
}

// processBatchPutRequest and return a batchPutResponse and error if encountered.
// This implementation aims to reduce database actions by sacrificing memory usage.
func (v *APIServerV1) processBatchPutRequest(r *http.Request, req batchPutRequest) (batchPutResponse, error) {
	resp := batchPutResponse{
		response: newSuccessResponse(),
	}

	txDb, tx, err := v.db.AsTx(r.Context(), nil)
	if err != nil {
		return resp, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Collect all class params into one slice.
	// Insert classes.
	// - For each created class, update its class groups' class_id foreign key.
	// - At the same time, collect all class group params into one slice.
	// Then insert the class groups.
	// - For each class group, update its class group sessions' class_group_id foreign key.
	// - At the same time, collect all class group session params into one slice.
	// Then insert the class group sessions.
	// Then insert the students.
	// Then for each of the sessions, insert an enrollment for each student.
	var (
		classesParams     []database.UpsertClassParams      // Store class params.
		classGroupsParams []database.UpsertClassGroupParams // Store class group params.

		// upsertClassGroups is a helper for class group session processing. The array is in the same order in which each
		// class group is created. This allows us to access and set variables within the class group using values
		// that are available to us only during creation of each class group.
		classGroups []*common.ClassGroupData

		classGroupSessionsParams []database.UpsertClassGroupSessionParams // Store class group session params.

		// users is a helper for session enrollment processing. This is a two-dimensional array, and is simply
		// an array of lists of students corresponding to each class group session.
		users [][]string

		usersParams       []database.UpsertUserParams              // Store user params.
		enrollmentsParams []database.UpsertSessionEnrollmentParams // Store session enrollment params.
	)

	for _, class := range req.Batches {
		classesParams = append(classesParams, class.Class)
	}

	upsertClasses, err := txDb.BatchUpsertClasses(r.Context(), classesParams)
	if err != nil {
		return resp, err
	}

	for i, class := range upsertClasses {
		resp.ClassIDs = append(resp.ClassIDs, class.ID)

		classData := &req.Batches[i]
		for idx := range classData.ClassGroups {
			classData.ClassGroups[idx].UpsertClassGroupParams.ClassID = class.ID
			classGroupsParams = append(classGroupsParams, classData.ClassGroups[idx].UpsertClassGroupParams)
			classGroups = append(classGroups, &classData.ClassGroups[idx])
		}
	}

	upsertClassGroups, err := txDb.BatchUpsertClassGroups(r.Context(), classGroupsParams)
	if err != nil {
		return resp, err
	}

	for i, group := range upsertClassGroups {
		classGroup := classGroups[i]
		usersParams = append(usersParams, classGroup.Students...)
		userIds := make([]string, 0, len(classGroup.Students))
		for _, user := range classGroup.Students {
			userIds = append(userIds, user.ID)
		}

		for idx := range classGroup.Sessions {
			classGroup.Sessions[idx].ClassGroupID = group.ID
			classGroupSessionsParams = append(classGroupSessionsParams, classGroup.Sessions[idx])
			users = append(users, userIds)
		}
	}

	upsertClassGroupSessions, err := txDb.BatchUpsertClassGroupSessions(r.Context(), classGroupSessionsParams)
	if err != nil {
		return resp, err
	}

	for i, session := range upsertClassGroupSessions {
		for idx := range users[i] {
			enrollmentsParams = append(enrollmentsParams, database.UpsertSessionEnrollmentParams{
				SessionID: session.ID,
				UserID:    users[i][idx],
			})
		}
	}

	if _, err = txDb.BatchUpsertUsers(r.Context(), usersParams); err != nil {
		return resp, err
	}

	if _, err = txDb.BatchUpsertSessionEnrollments(r.Context(), enrollmentsParams); err != nil {
		return resp, err
	}

	return resp, tx.Commit()
}
