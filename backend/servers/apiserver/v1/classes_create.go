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
	Classes []common.ClassCreationData `json:"classes"`
}

// classesCreate is the handler for a request to create classes.
func (v *APIServerV1) classesCreate(w http.ResponseWriter, r *http.Request) {
	var (
		resp apiResponse
		err  error
	)

	switch contentType := r.Header.Get("Content-Type"); {
	case strings.HasPrefix(contentType, "multipart"):
		resp, err = v.processClassCreationFiles(r)
	case contentType == "application/json":
		resp, err = v.processClassCreationJSON(r)
	default:
		resp = newErrorResponse(http.StatusUnsupportedMediaType, fmt.Sprintf("%s is unsupported", contentType))
	}

	if err != nil {
		resp = newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	v.writeResponse(w, classesUrl, resp)
}

// processClassCreationFiles processes a request to create classes via file uploads.
func (v *APIServerV1) processClassCreationFiles(r *http.Request) (apiResponse, error) {
	var resp apiResponse

	if err := r.ParseMultipartForm(maxParseMemory); err != nil {
		return resp, err
	}

	limiter := goroutines.NewLimiter(maxGoRoutines)

	saveRes := sync.Map{}
	for _, header := range r.MultipartForm.File[multipartFormFileIdent] {
		header := header // Required for go routine to point to different file for each loop.
		limiter.Do(func() {
			creationData, err := v.processClassCreationFile(header)
			saveRes.Store(&creationData, err)
		})
	}

	limiter.Wait()

	var (
		request classesCreateRequest
		err     error
	)
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

		request.Classes = append(request.Classes, *data)
		return true
	})

	if err != nil {
		return resp, err
	}

	if resp, err = v.processClassesCreateRequest(r.Context(), request); err != nil {
		return resp, err
	}

	return resp, nil
}

// processClassCreationFile processes a file to create a new class.
func (v *APIServerV1) processClassCreationFile(fileHeader *multipart.FileHeader) (common.ClassCreationData, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return common.ClassCreationData{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	v.l.Debug(fmt.Sprintf("%s - processing class creation file", namespace),
		zap.String("filename", fileHeader.Filename))

	creationData, err := common.ParseClassCreationFile(fileHeader.Filename, file)
	if err != nil {
		return creationData, fmt.Errorf("%s - error parsing class creation file %s: %w", namespace, fileHeader.Filename, err)
	}

	return creationData, nil
}

// processClassCreationJSON processes a request to create classes via JSON body.
func (v *APIServerV1) processClassCreationJSON(r *http.Request) (apiResponse, error) {
	var (
		resp apiResponse
		b    bytes.Buffer
	)

	if _, err := b.ReadFrom(r.Body); err != nil {
		return resp, err
	}

	var request classesCreateRequest
	if err := json.Unmarshal(b.Bytes(), &request); err != nil {
		return resp, err
	}

	resp, err := v.processClassesCreateRequest(r.Context(), request)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// processClassesCreateRequest and return a classesCreateResponse and error if encountered.
func (v *APIServerV1) processClassesCreateRequest(ctx context.Context, request classesCreateRequest) (apiResponse, error) {
	var resp apiResponse

	if err := request.isValid(); err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error()), nil
	}

	tx, err := v.db.C.Begin(ctx)
	if err != nil {
		return resp, err
	}

	q := v.db.Q.WithTx(tx)

	var coursesParams []database.UpsertCoursesParams
	for _, class := range request.Classes {
		coursesParams = append(coursesParams, class.Course)
	}

	q.UpsertCourses(ctx, coursesParams).QueryRow(v.upsertClassGroups(ctx, q, &resp))

	if !resp.wasSuccessful() {
		return resp, tx.Rollback(ctx)
	} else if err = tx.Commit(ctx); err != nil {
		return resp, errors.Join(err, tx.Rollback(ctx))
	}

	return resp, nil
}

func (v *APIServerV1) upsertClassGroups(ctx context.Context, q *database.Queries, resp *classesCreateResponse) func(i int, course database.Course, err error) {
	return func(i int, course database.Course, err error) {
		if !resp.wasSuccessful() {
			return
		}

		class := &resp.Classes[i]

		if err != nil {
			class.Error = err.Error()
			return
		}

		var classGroupsParams []database.UpsertClassGroupsParams
		for idx := range class.ClassGroups {
			class.ClassGroups[idx].UpsertClassGroupsParams.CourseID = course.ID
			classGroupsParams = append(classGroupsParams, class.ClassGroups[idx].UpsertClassGroupsParams)
		}

		q.UpsertClassGroups(ctx, classGroupsParams).QueryRow(v.upsertSessionsAndUsers(ctx, q, resp, class))
	}
}

func (v *APIServerV1) upsertSessionsAndUsers(
	ctx context.Context,
	q *database.Queries,
	resp *classesCreateResponse,
	class *common.ClassCreationData,
) func(i int, group database.ClassGroup, err error) {
	return func(i int, group database.ClassGroup, err error) {
		if !resp.wasSuccessful() {
			return
		}

		classGroup := class.ClassGroups[i] // Do not take pointer to avoid editing.

		if err != nil {
			class.Error = err.Error()
			return
		}

		for idx := range classGroup.Sessions {
			session := &classGroup.Sessions[idx]
			session.ClassGroupID = group.ID
		}

		// Insert the sessions for this class group.
		q.UpsertClassGroupSessions(ctx, classGroup.Sessions).QueryRow(func(i int, session database.ClassGroupSession, err error) {
			if err != nil {
				(&resp.Classes[i]).Error = err.Error()
				return
			}
		})
	}
}
