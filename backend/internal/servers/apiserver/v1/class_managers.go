package v1

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
)

const (
	maxClassManagersPutParseMemory      = 32 << 20
	maxClassManagersPutFiles            = 1
	multipartFormClassManagersFileIdent = "manager-attachments"
)

func (v *APIServerV1) classManagers(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.classManagersGet(r)
	case http.MethodPost:
		resp = v.classManagersPost(r)
	case http.MethodPut:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classManagersGetResponse struct {
	response
	ClassManagers []model.ClassManager `json:"class_managers"`
}

func (v *APIServerV1) classManagersGet(r *http.Request) apiResponse {
	params, err := database.DecodeListQueryParams(
		r.URL.Query(), table.ClassManagers, table.ClassManagers.AllColumns,
	)
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}

	managers, err := v.db.ListClassManagers(r.Context(), params)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class managers get database action")
	}

	resp := classManagersGetResponse{
		newSuccessResponse(),
		make([]model.ClassManager, 0, len(managers)),
	}

	resp.ClassManagers = append(resp.ClassManagers, managers...)
	return resp
}

type classManagersPostRequest struct {
	ClassManager database.CreateClassManagerParams `json:"class_manager"`
}

type classManagersPostResponse struct {
	response
	ClassManager classManagersPostClassManagerResponseFields `json:"class_manager"`
}

type classManagersPostClassManagerResponseFields struct {
	ID           int64              `json:"id"`
	UserID       string             `json:"user_id"`
	ClassID      int64              `json:"class_id"`
	ManagingRole model.ManagingRole `json:"managing_role"`
	CreatedAt    time.Time          `json:"created_at"`
}

func (v *APIServerV1) classManagersPost(r *http.Request) apiResponse {
	var req classManagersPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	manager, err := v.db.CreateClassManager(r.Context(), req.ClassManager)
	if err != nil {
		switch {
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusConflict, "class manager with same user_id and class_id already exists")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusBadRequest, "user_id and/or class_id does not exist")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process class managers post database action")
		}
	}

	return classManagersPostResponse{
		newSuccessResponse(),
		classManagersPostClassManagerResponseFields{
			manager.ID,
			manager.UserID,
			manager.ClassID,
			manager.ManagingRole,
			manager.CreatedAt,
		},
	}
}

type classManagersPutResponse struct {
	response
}

func (v *APIServerV1) classManagerPut(r *http.Request) apiResponse {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart") {
		return newErrorResponse(http.StatusUnsupportedMediaType, "a multipart request body is required")
	}

	resp, err := v.processClassManagerPutRequest(r)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class manager put file(s)")
	}

	return resp
}

func (v *APIServerV1) processClassManagerPutRequest(r *http.Request) (apiResponse, error) {
	if err := r.ParseMultipartForm(maxClassManagersPutParseMemory); err != nil {
		return classManagersPutResponse{}, err
	}

	if len(r.MultipartForm.File[multipartFormClassManagersFileIdent]) != maxClassManagersPutFiles {
		return newErrorResponse(http.StatusBadRequest, "only one file is allowed"), nil
	}

	// TODO: Implement class managers file parsing.
	return newErrorResponse(http.StatusNotImplemented, ""), nil
}
