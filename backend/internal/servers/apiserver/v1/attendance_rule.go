package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) attendanceRule(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, attendanceRuleUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.attendanceRuleGet(r, classId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type attendanceRuleGetResponse struct {
	response
	CoordinatingClass database.CoordinatingClass  `json:"coordinating_class"`
	Rules             []model.ClassAttendanceRule `json:"rules"`
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) attendanceRuleGet(r *http.Request, id int64) apiResponse {
	class, err := v.db.GetCoordinatingClass(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the request coordinating class does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process attendance rule get database action")
	}

	rules, err := v.db.GetCoordinatingClassRules(r.Context(), id)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get coordinating class rules")
	}

	return attendanceRuleGetResponse{
		newSuccessResponse(),
		class,
		append(make([]model.ClassAttendanceRule, 0, len(rules)), rules...),
	}
}
