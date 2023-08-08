package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/common"
	"github.com/darylhjd/oams/backend/pkg/goroutines"
)

const (
	maxParseMemory         = 32 << 20
	maxGoRoutines          = 10
	multipartFormFileIdent = "attachments"
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

	resp, err := v.processPostBody(r)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process batch file(s)")
	}

	return resp
}

func (v *APIServerV1) processPostBody(r *http.Request) (apiResponse, error) {
	if err := r.ParseMultipartForm(maxParseMemory); err != nil {
		return batchPostResponse{}, err
	}

	limiter := goroutines.NewLimiter(maxGoRoutines)

	saveRes := sync.Map{}
	for _, header := range r.MultipartForm.File[multipartFormFileIdent] {
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

			data, err = common.ParseBatchFile(header.Filename, file)
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
		err           error
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
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	resp, err := v.processBatchPutRequest(r, req)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return resp
}

// processBatchPutRequest and return a batchPutResponse and error if encountered.
// This implementation aims to reduce database actions by sacrificing memory usage.
func (v *APIServerV1) processBatchPutRequest(r *http.Request, req batchPutRequest) (batchPutResponse, error) {
	resp := batchPutResponse{
		response: newSuccessResponse(),
	}

	var dbErr error
	tx, err := v.db.C.Begin(r.Context())
	if err != nil {
		return resp, err
	}
	defer func() {
		_ = tx.Rollback(r.Context())
	}()

	q := v.db.Q.WithTx(tx)

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
		classesParams     []database.UpsertClassesParams     // Store class params.
		classGroupsParams []database.UpsertClassGroupsParams // Store class group params.

		// classGroups is a helper for class group session processing. The array is in the same order in which each
		// class group is created. This allows us to access and set variables within the class group using values
		// that are available to us only during creation of each class group.
		classGroups []*common.ClassGroupData

		classGroupSessionsParams []database.UpsertClassGroupSessionsParams // Store class group session params.

		// users is a helper for session enrollment processing. This is a two-dimensional array, and is simply
		// an array of lists of students corresponding to each class group session.
		users [][]string

		usersParams       []database.UpsertUsersParams              // Store user params.
		enrollmentsParams []database.UpsertSessionEnrollmentsParams // Store session enrollment params/
	)

	for _, class := range req.Batches {
		classesParams = append(classesParams, class.Class)
	}

	q.UpsertClasses(r.Context(), classesParams).QueryRow(func(i int, class database.Class, err error) {
		if dbErr != nil {
			return
		} else if err != nil {
			dbErr = err
			return
		}

		resp.ClassIDs = append(resp.ClassIDs, class.ID)

		classData := &req.Batches[i]
		for idx := range classData.ClassGroups {
			classData.ClassGroups[idx].UpsertClassGroupsParams.ClassID = class.ID
			classGroupsParams = append(classGroupsParams, classData.ClassGroups[idx].UpsertClassGroupsParams)
			classGroups = append(classGroups, &classData.ClassGroups[idx])
		}
	})
	if dbErr != nil {
		return resp, dbErr
	}

	q.UpsertClassGroups(r.Context(), classGroupsParams).QueryRow(func(i int, group database.ClassGroup, err error) {
		if dbErr != nil {
			return
		} else if err != nil {
			dbErr = err
			return
		}

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
	})
	if dbErr != nil {
		return resp, dbErr
	}

	q.UpsertClassGroupSessions(r.Context(), classGroupSessionsParams).QueryRow(func(i int, session database.ClassGroupSession, err error) {
		if dbErr != nil {
			return
		} else if err != nil {
			dbErr = err
			return
		}

		for idx := range users[i] {
			enrollmentsParams = append(enrollmentsParams, database.UpsertSessionEnrollmentsParams{
				SessionID: session.ID,
				UserID:    users[i][idx],
			})
		}
	})
	if dbErr != nil {
		return resp, dbErr
	}

	if dbErr = q.UpsertUsers(r.Context(), usersParams).Close(); dbErr != nil {
		return resp, dbErr
	}

	if dbErr = q.UpsertSessionEnrollments(r.Context(), enrollmentsParams).Close(); dbErr != nil {
		return resp, dbErr
	}

	return resp, tx.Commit(r.Context())
}
