package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/common"
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

func (v *APIServerV1) batchPost(r *http.Request) apiResponse {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart") {
		return newErrorResponse(http.StatusUnsupportedMediaType, "a multipart request body is required")
	}

	req, err := v.fromBatchFiles(r)
	if err != nil {
	}

	return batchPostResponse{
		newSuccessResponse(),
		req,
	}
}

// fromBatchFiles creates a request struct from uploaded files.
func (v *APIServerV1) fromBatchFiles(r *http.Request) (batchPutRequest, error) {
	var req batchPutRequest

	if err := r.ParseMultipartForm(maxParseMemory); err != nil {
		return req, err
	}

	limiter := goroutines.NewLimiter(maxGoRoutines)

	saveRes := sync.Map{}
	for _, header := range r.MultipartForm.File[multipartFormFileIdent] {
		header := header // Required for go routine to point to different file for each loop.
		limiter.Do(func() {
			creationData, err := v.fromBatchFile(header)
			saveRes.Store(&creationData, err)
		})
	}

	limiter.Wait()

	var err error
	saveRes.Range(func(key, value any) bool {
		data, ok := key.(*common.BatchData)
		if !ok {
			err = errors.New("type assertion failed when processing class creation data")
			return false
		}

		if value != nil {
			err = value.(error)
			return false
		}

		req.Batches = append(req.Batches, *data)
		return true
	})

	return req, err
}

// fromBatchFile processes a file to create new class creation data.
func (v *APIServerV1) fromBatchFile(fileHeader *multipart.FileHeader) (common.BatchData, error) {
	var data common.BatchData

	file, err := fileHeader.Open()
	if err != nil {
		return data, err
	}
	defer func() {
		_ = file.Close()
	}()

	v.l.Debug(fmt.Sprintf("%s - processing class creation file", namespace), zap.String("filename", fileHeader.Filename))

	data, err = common.ParseBatchFile(fileHeader.Filename, file)
	if err != nil {
		return data, err
	}

	return data, nil
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

// batchPut is the handler for a request to create a batch of entities.
func (v *APIServerV1) batchPut(r *http.Request) apiResponse {
	var (
		req  batchPutRequest
		resp apiResponse
		err  error
	)

	switch contentType := r.Header.Get("Content-Type"); {
	case strings.HasPrefix(contentType, "multipart"):
		req, err = v.fromBatchFiles(r)
	case contentType == "application/json":
		req, err = v.fromBatchJSON(r)
	default:
		return newErrorResponse(http.StatusUnsupportedMediaType, fmt.Sprintf("%s is unsupported", contentType))
	}

	if err == nil {
		resp, err = v.processBatchPutRequest(r.Context(), req)
	}

	if err != nil {
		v.logInternalServerError(r, err)
		resp = newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return resp
}

// fromBatchJSON creates a request struct from JSON body.
func (v *APIServerV1) fromBatchJSON(r *http.Request) (batchPutRequest, error) {
	var (
		req batchPutRequest
		b   bytes.Buffer
	)

	if _, err := b.ReadFrom(r.Body); err != nil {
		return req, err
	}

	err := json.Unmarshal(b.Bytes(), &req)
	return req, err
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
func (v *APIServerV1) processBatchPutRequest(ctx context.Context, req batchPutRequest) (apiResponse, error) {
	tx, err := v.db.C.Begin(ctx)
	if err != nil {
		return nil, err
	}

	q := v.db.Q.WithTx(tx)

	var (
		dbErr         error
		coursesParams []database.UpsertClassesParams
	)
	resp := batchPutResponse{
		response: newSuccessResponse(),
	}

	defer func() {
		_ = tx.Rollback(ctx)
		if dbErr != nil {
			v.l.Debug(fmt.Sprintf("%s - error while doing classes create action", namespace), zap.Error(dbErr))
		}
	}()

	for _, class := range req.Batches {
		coursesParams = append(coursesParams, class.Course)
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
		return nil, dbErr
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
		return nil, dbErr
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
				session.ID,
				classGroupSessionsHelper.students[i][idx].ID,
				false,
			})
		}
	})
	if dbErr != nil {
		return nil, dbErr
	}

	students, dbErr := upsertUsers(ctx, v.db, tx, studentsParams)
	if dbErr != nil {
		return nil, dbErr
	}

	resp.Students = len(students)

	if dbErr = q.CreateSessionEnrollments(ctx, enrollmentsParams).Close(); dbErr != nil {
		return nil, dbErr
	}

	sessionEnrollments, dbErr := q.ListSessionEnrollments(ctx)
	if dbErr != nil {
		return nil, dbErr
	}

	resp.SessionEnrollments = len(sessionEnrollments)
	return resp, tx.Commit(ctx)
}
