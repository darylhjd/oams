package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/common"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

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
	Classes            int `json:"classes"`
	ClassGroups        int `json:"class_groups"`
	ClassGroupSessions int `json:"class_group_sessions"`
	Students           int `json:"students"`
	SessionEnrollments int `json:"session_enrollments"`
}

// batchPut is the handler that does the actual processing of the entities.
func (v *APIServerV1) batchPut(r *http.Request) apiResponse {
	var req batchPutRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	resp, err := v.processBatchPutRequest(r.Context(), req)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return resp
}

type classGroupsParamsWithClassGroup struct {
	classGroupsParams []database.UpsertClassGroupsParams
	classGroups       []*common.ClassGroupData
}

type classGroupSessionsParamsWithStudents struct {
	classGroupSessionsParams []database.UpsertClassGroupSessionsParams
	students                 [][]database.UpsertUsersParams
}

// processBatchPutRequest and return a batchPutResponse and error if encountered.
func (v *APIServerV1) processBatchPutRequest(ctx context.Context, req batchPutRequest) (batchPutResponse, error) {
	resp := batchPutResponse{
		response: newSuccessResponse(),
	}

	var (
		dbErr         error
		coursesParams []database.UpsertClassesParams
	)
	tx, err := v.db.C.Begin(ctx)
	if err != nil {
		return resp, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
		if dbErr != nil {
			v.l.Debug(fmt.Sprintf("%s - error while doing classes create action", namespace), zap.Error(dbErr))
		}
	}()

	q := v.db.Q.WithTx(tx)

	for _, class := range req.Batches {
		coursesParams = append(coursesParams, class.Class)
	}

	// Insert courses.
	// Then, for each course, update its class group's course_id foreign key.
	// Then insert the class groups.
	// Then for each class group, for each of its sessions, update its class_group_id.
	// Then insert the sessions.
	// Then insert the students.
	// Then for each of the sessions, insert an enrollment for each student.
	var classGroupsHelper classGroupsParamsWithClassGroup
	q.UpsertClasses(ctx, coursesParams).QueryRow(func(i int, course database.Class, err error) {
		if dbErr != nil {
			return
		} else if err != nil {
			dbErr = err
			return
		}

		resp.Classes++

		class := req.Batches[i]
		for idx := range class.ClassGroups {
			class.ClassGroups[idx].UpsertClassGroupsParams.ClassID = course.ID
			classGroupsHelper.classGroupsParams = append(classGroupsHelper.classGroupsParams, class.ClassGroups[idx].UpsertClassGroupsParams)
			classGroupsHelper.classGroups = append(classGroupsHelper.classGroups, &class.ClassGroups[idx])
		}
	})
	if dbErr != nil {
		return resp, dbErr
	}

	var classGroupSessionsHelper classGroupSessionsParamsWithStudents
	q.UpsertClassGroups(ctx, classGroupsHelper.classGroupsParams).QueryRow(func(i int, group database.ClassGroup, err error) {
		if dbErr != nil {
			return
		} else if err != nil {
			dbErr = err
			return
		}

		resp.ClassGroups++

		for idx := range classGroupsHelper.classGroups[i].Sessions {
			classGroupsHelper.classGroups[i].Sessions[idx].ClassGroupID = group.ID
			classGroupSessionsHelper.classGroupSessionsParams = append(classGroupSessionsHelper.classGroupSessionsParams, classGroupsHelper.classGroups[i].Sessions[idx].UpsertClassGroupSessionsParams)
			classGroupSessionsHelper.students = append(classGroupSessionsHelper.students, classGroupsHelper.classGroups[i].Students)
		}
	})
	if dbErr != nil {
		return resp, dbErr
	}

	var (
		studentsParams    []database.UpsertUsersParams
		enrollmentsParams []database.CreateSessionEnrollmentsParams
	)
	q.UpsertClassGroupSessions(ctx, classGroupSessionsHelper.classGroupSessionsParams).QueryRow(func(i int, session database.ClassGroupSession, err error) {
		if dbErr != nil {
			return
		} else if err != nil {
			dbErr = err
			return
		}

		resp.ClassGroupSessions++

		for idx := range classGroupSessionsHelper.students[i] {
			studentsParams = append(studentsParams, classGroupSessionsHelper.students[i][idx])
			enrollmentsParams = append(enrollmentsParams, database.CreateSessionEnrollmentsParams{
				SessionID: session.ID,
				UserID:    classGroupSessionsHelper.students[i][idx].ID,
			})
		}
	})
	if dbErr != nil {
		return resp, dbErr
	}

	students, dbErr := upsertUsers(ctx, v.db, tx, studentsParams)
	if dbErr != nil {
		return resp, dbErr
	}

	resp.Students = len(students)

	if dbErr = q.CreateSessionEnrollments(ctx, enrollmentsParams).Close(); dbErr != nil {
		return resp, dbErr
	}

	sessionEnrollments, dbErr := q.ListSessionEnrollments(ctx)
	if dbErr != nil {
		return resp, dbErr
	}

	resp.SessionEnrollments = len(sessionEnrollments)
	return resp, tx.Commit(ctx)
}

// upsertUsers inserts the provided usersParams into the specified database. If tx is nil, a new transaction is started.
// Otherwise, a nested transaction (using save points) is used.
func upsertUsers(ctx context.Context, db *database.DB, tx pgx.Tx, usersParams []database.UpsertUsersParams) ([]database.User, error) {
	var err error

	if tx != nil {
		tx, err = tx.Begin(ctx)
	} else {
		tx, err = db.C.Begin(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	q := db.Q.WithTx(tx)

	if err = q.UpsertUsers(ctx, usersParams).Close(); err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(usersParams))
	for _, param := range usersParams {
		ids = append(ids, param.ID)
	}

	users, err := q.GetUsersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return users, tx.Commit(ctx)
}
