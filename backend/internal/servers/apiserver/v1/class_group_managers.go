package v1

import (
	"net/http"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/common"
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
		resp = newErrorResponse(http.StatusNotImplemented, "")
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

	multipartFile := r.MultipartForm.File[multipartFormClassGroupManagersFileIdent][maxClassGroupManagersPutFiles-1]
	file, err := multipartFile.Open()
	if err != nil {
		return nil, err
	}

	upsertData, err := common.ParseManagersFile(file)
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error()), nil
	}

	_, err = v.db.BatchUpsertClassGroupManagers(r.Context(), upsertData)
	if err != nil {
		return nil, err
	}

	return newErrorResponse(http.StatusNotImplemented, ""), nil
}
