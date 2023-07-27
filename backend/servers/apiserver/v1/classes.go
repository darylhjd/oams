package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) classes(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.classesGet(r)
	case http.MethodPost:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, classesUrl, resp)
}

type classesGetResponse struct {
	response
	Classes []database.Class `json:"classes"`
}

func (v *APIServerV1) classesGet(r *http.Request) apiResponse {
	classes, err := v.db.Q.ListClasses(r.Context())
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, "could not process classes get database action")
	}

	resp := classesGetResponse{
		newSuccessResponse(),
		make([]database.Class, 0, len(classes)),
	}

	resp.Classes = append(resp.Classes, classes...)
	return resp
}
