package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) attendanceRules(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.attendanceRulesGet(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type attendanceRulesGetResponse struct {
	response
	CoordinatingClasses []database.CoordinatingClasses `json:"coordinating_classes"`
}

func (v *APIServerV1) attendanceRulesGet(r *http.Request) apiResponse {
	classes, err := v.db.GetCoordinatingClasses(r.Context())
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process attendance rules get database action")
	}

	return attendanceRulesGetResponse{
		newSuccessResponse(),
		append(make([]database.CoordinatingClasses, 0, len(classes)), classes...),
	}
}
