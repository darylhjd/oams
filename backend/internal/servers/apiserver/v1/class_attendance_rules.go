package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
)

func (v *APIServerV1) classAttendanceRules(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.classAttendanceRulesGet(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classAttendanceRulesGetResponse struct {
	response
	ClassAttendanceRules []model.ClassAttendanceRule `json:"class_attendance_rules"`
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) classAttendanceRulesGet(r *http.Request) apiResponse {
	params, err := database.DecodeListQueryParams(r.URL.Query(), table.ClassAttendanceRules.AllColumns)
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}

	rules, err := v.db.ListClassAttendanceRules(r.Context(), params)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class attendance rules get database action")
	}

	return classAttendanceRulesGetResponse{
		newSuccessResponse(),
		append(make([]model.ClassAttendanceRule, 0, len(rules)), rules...),
	}
}
