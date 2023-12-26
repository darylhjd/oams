package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
)

func (v *APIServerV1) classGroups(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupsGet(r)
	case http.MethodPost:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGroupsGetResponse struct {
	response
	ClassGroups []model.ClassGroup `json:"class_groups"`
}

func (v *APIServerV1) classGroupsGet(r *http.Request) apiResponse {
	params, err := database.DecodeListQueryParams(r.URL.Query(), table.ClassGroups.AllColumns)
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}

	groups, err := v.db.ListClassGroups(r.Context(), params)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class groups get database action")
	}

	resp := classGroupsGetResponse{
		newSuccessResponse(),
		make([]model.ClassGroup, 0, len(groups)),
	}

	resp.ClassGroups = append(resp.ClassGroups, groups...)
	return resp
}
