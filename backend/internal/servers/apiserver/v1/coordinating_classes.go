package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) coordinatingClasses(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.coordinatingClassesGet(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type coordinatingClassesGetResponse struct {
	response
	CoordinatingClasses []database.CoordinatingClass `json:"coordinating_classes"`
}

func (v *APIServerV1) coordinatingClassesGet(r *http.Request) apiResponse {
	classes, err := v.db.GetCoordinatingClasses(r.Context())
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process coordinating classes get database action")
	}

	return coordinatingClassesGetResponse{
		newSuccessResponse(),
		append(make([]database.CoordinatingClass, 0, len(classes)), classes...),
	}
}
