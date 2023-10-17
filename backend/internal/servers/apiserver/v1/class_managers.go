package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
)

func (v *APIServerV1) classManagers(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.classManagersGet(r)
	case http.MethodPost:
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
	params, err := v.decodeListQueryParameters(r.URL.Query(), table.ClassManagers.AllColumns)
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
