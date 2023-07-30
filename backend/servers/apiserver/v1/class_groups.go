package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) classGroups(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodPost:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, classGroupsUrl, resp)
}

type classGroupsGetResponse struct {
	response
	ClassGroups []database.ClassGroup `json:"class_groups"`
}

func (v *APIServerV1) classGroupsGet(r *http.Request) apiResponse {
	groups, err := v.db.Q.ListClassGroups(r.Context())
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, "could not process class groups get database action")
	}

	resp := classGroupsGetResponse{
		newSuccessResponse(),
		make([]database.ClassGroup, 0, len(groups)),
	}

	resp.ClassGroups = append(resp.ClassGroups, groups...)
	return resp
}
