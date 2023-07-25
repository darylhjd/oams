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
	"github.com/darylhjd/oams/backend/servers/apiserver/common"
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/pkg/goroutines"
)

const (
	maxParseMemory         = 32 << 20
	maxGoRoutines          = 10
	multipartFormFileIdent = "attachments"
)

type classesCreateRequest struct {
	Classes []common.ClassCreationData `json:"classes"`
}

// isValid does a validation of the request, and returns an error if it is not valid.
func (r classesCreateRequest) isValid() error {
	for _, class := range r.Classes {
		if err := class.IsValid(); err != nil {
			return err
		}
	}

	return nil
}

type classesCreateResponse struct {
	response
	Classes            int `json:"classes"`
	ClassGroups        int `json:"class_groups"`
	ClassGroupSessions int `json:"class_group_sessions"`
	Students           int `json:"students"`
	SessionEnrollments int `json:"session_enrollments"`
}

// classesCreate is the handler for a request to create classes.
func (v *APIServerV1) classesCreate(w http.ResponseWriter, r *http.Request) {
	var (
		req  classesCreateRequest
		resp apiResponse
		err  error
	)

	switch contentType := r.Header.Get("Content-Type"); {
	case strings.HasPrefix(contentType, "multipart"):
		req, err = v.fromClassCreationFiles(r)
	case contentType == "application/json":
		req, err = v.fromClassCreationJSON(r)
	default:
		resp = newErrorResponse(http.StatusUnsupportedMediaType, fmt.Sprintf("%s is unsupported", contentType))
		v.writeResponse(w, classesUrl, resp)
		return
	}

	if err == nil {
		resp, err = v.processClassesCreateRequest(r.Context(), req)
	}

	if err != nil {
		resp = newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	v.writeResponse(w, classesUrl, resp)
}

// fromClassCreationFiles creates a request struct from uploaded files.
func (v *APIServerV1) fromClassCreationFiles(r *http.Request) (classesCreateRequest, error) {
	var req classesCreateRequest

	if err := r.ParseMultipartForm(maxParseMemory); err != nil {
		return req, err
	}

	limiter := goroutines.NewLimiter(maxGoRoutines)

	saveRes := sync.Map{}
	for _, header := range r.MultipartForm.File[multipartFormFileIdent] {
		header := header // Required for go routine to point to different file for each loop.
		limiter.Do(func() {
			creationData, err := v.fromClassCreationFile(header)
			saveRes.Store(&creationData, err)
		})
	}

	limiter.Wait()

	var err error
	saveRes.Range(func(key, value any) bool {
		data, ok := key.(*common.ClassCreationData)
		if !ok {
			err = errors.New("type assertion failed when processing class creation data")
			return false
		}

		if value != nil {
			err = value.(error)
			return false
		}

		req.Classes = append(req.Classes, *data)
		return true
	})

	return req, err
}

// fromClassCreationFile processes a file to create new class creation data.
func (v *APIServerV1) fromClassCreationFile(fileHeader *multipart.FileHeader) (common.ClassCreationData, error) {
	var data common.ClassCreationData

	file, err := fileHeader.Open()
	if err != nil {
		return data, err
	}
	defer func() {
		_ = file.Close()
	}()

	v.l.Debug(fmt.Sprintf("%s - processing class creation file", namespace),
		zap.String("filename", fileHeader.Filename))

	data, err = common.ParseClassCreationFile(fileHeader.Filename, file)
	if err != nil {
		return data, fmt.Errorf("%s - error parsing class creation file %s: %w", namespace, fileHeader.Filename, err)
	}

	return data, nil
}

// fromClassCreationJSON creates a request struct from JSON body.
func (v *APIServerV1) fromClassCreationJSON(r *http.Request) (classesCreateRequest, error) {
	var (
		req classesCreateRequest
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

// processClassesCreateRequest and return a classesCreateResponse and error if encountered.
func (v *APIServerV1) processClassesCreateRequest(ctx context.Context, req classesCreateRequest) (apiResponse, error) {
	if err := req.isValid(); err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error()), nil
	}

	tx, err := v.db.C.Begin(ctx)
	if err != nil {
		return nil, err
	}

	q := v.db.Q.WithTx(tx)

	var (
		dbErr         error
		coursesParams []database.UpsertClassesParams
	)
	resp := classesCreateResponse{
		response: newSuccessfulResponse(),
	}

	defer func() {
		_ = tx.Rollback(ctx)
		if dbErr != nil {
			v.l.Debug(fmt.Sprintf("%s - error while doing classes create action", namespace),
				zap.Error(dbErr))
		}
	}()

	for _, class := range req.Classes {
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

		class := req.Classes[i]
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
				SessionID: session.ID,
				UserID:    classGroupSessionsHelper.students[i][idx].ID,
			})
		}
	})
	if dbErr != nil {
		return nil, dbErr
	}

	if dbErr = q.UpsertUsers(ctx, studentsParams).Close(); dbErr != nil {
		return nil, dbErr
	}

	if dbErr = q.CreateSessionEnrollments(ctx, enrollmentsParams).Close(); dbErr != nil {
		return nil, dbErr
	}

	students, dbErr := q.ListUsers(ctx)
	if dbErr != nil {
		return nil, dbErr
	}

	resp.Students = len(students)

	sessionEnrollments, dbErr := q.ListSessionEnrollments(ctx)
	if dbErr != nil {
		return nil, dbErr
	}

	resp.SessionEnrollments = len(sessionEnrollments)
	return resp, tx.Commit(ctx)
}
