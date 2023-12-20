package v1

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/common/managers"
)

const (
	maxClassGroupManagersPutParseMemory      = 32 << 20
	maxClassGroupManagersPutFiles            = 1
	multipartFormClassGroupManagersFileIdent = "manager-attachments"
)

func (v *APIServerV1) classGroupManagers(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupManagersGet(r)
	case http.MethodPost:
		resp = v.classGroupManagersPost(r)
	case http.MethodPut:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGroupManagersGetResponse struct {
	response
	ClassGroupManagers []model.ClassGroupManager `json:"class_group_managers"`
}

func (v *APIServerV1) classGroupManagersGet(r *http.Request) apiResponse {
	params, err := database.DecodeListQueryParams(
		r.URL.Query(), table.ClassGroupManagers, table.ClassGroupManagers.AllColumns,
	)
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}

	managers, err := v.db.ListClassGroupManagers(r.Context(), params)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group managers get database action")
	}

	resp := classGroupManagersGetResponse{
		newSuccessResponse(),
		make([]model.ClassGroupManager, 0, len(managers)),
	}

	resp.ClassGroupManagers = append(resp.ClassGroupManagers, managers...)
	return resp
}

type classGroupManagersPostRequest struct {
	ClassGroupManager database.CreateClassGroupManagerParams `json:"class_group_manager"`
}

type classGroupManagersPostResponse struct {
	response
	ClassGroupManager classGroupManagersPostClassGroupManagerResponseFields `json:"class_group_manager"`
}

type classGroupManagersPostClassGroupManagerResponseFields struct {
	ID           int64              `json:"id"`
	UserID       string             `json:"user_id"`
	ClassGroupID int64              `json:"class_group_id"`
	ManagingRole model.ManagingRole `json:"managing_role"`
	CreatedAt    time.Time          `json:"created_at"`
}

func (v *APIServerV1) classGroupManagersPost(r *http.Request) apiResponse {
	var req classGroupManagersPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	manager, err := v.db.CreateClassGroupManager(r.Context(), req.ClassGroupManager)
	if err != nil {
		switch {
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusConflict, "class group manager with same user_id and class_group_id already exists")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusBadRequest, "user_id and/or class_group_id does not exist")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process class group managers post database action")
		}
	}

	return classGroupManagersPostResponse{
		newSuccessResponse(),
		classGroupManagersPostClassGroupManagerResponseFields{
			manager.ID,
			manager.UserID,
			manager.ClassGroupID,
			manager.ManagingRole,
			manager.CreatedAt,
		},
	}
}

type classGroupManagersPutResponse struct {
	response
}

func (v *APIServerV1) classGroupManagerPut(r *http.Request) apiResponse {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart") {
		return newErrorResponse(http.StatusUnsupportedMediaType, "a multipart request body is required")
	}

	resp, err := v.processClassGroupManagerPutRequest(r)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group manager put file(s)")
	}

	return resp
}

func (v *APIServerV1) processClassGroupManagerPutRequest(r *http.Request) (apiResponse, error) {
	if err := r.ParseMultipartForm(maxClassGroupManagersPutParseMemory); err != nil {
		return nil, err
	}

	if len(r.MultipartForm.File[multipartFormClassGroupManagersFileIdent]) != maxClassGroupManagersPutFiles {
		return newErrorResponse(http.StatusBadRequest, "only one file is allowed"), nil
	}

	// TODO: Implement class group managers file parsing.
	multipartFile := r.MultipartForm.File[multipartFormClassGroupManagersFileIdent][maxClassGroupManagersPutFiles-1]
	file, err := multipartFile.Open()
	if err != nil {
		return nil, err
	}

	upsertData, err := managers.ParseManagersFile(file)
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error()), nil
	}

	res, err := v.db.BatchUpsertClassGroupManagers(r.Context(), upsertData)
	if err != nil {
		return nil, err
	}

	return classGroupManagersPutResponse{newSuccessResponse()}, nil
}
