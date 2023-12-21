package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) attendanceTaking(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classGroupSessionId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, attendanceTakingUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group session id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.attendanceTakingGet(r, classGroupSessionId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type attendanceTakingGetResponse struct {
	response
	UpcomingClassGroupSession upcomingClassGroupSession `json:"upcoming_class_group_session"`
	EnrollmentData            []model.SessionEnrollment `json:"enrollment_data"`
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) attendanceTakingGet(r *http.Request, id int64) apiResponse {
	resp := attendanceTakingGetResponse{
		response: newSuccessResponse(),
	}

	upcoming, err := v.db.GetUpcomingManagedClassGroupSession(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested upcoming class group session does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process attendance taking get database action")
	}

	resp.UpcomingClassGroupSession = upcomingClassGroupSession{}.fromDatabaseUpcomingClassGroupSession(upcoming)
	enrollments, err := v.db.GetUpcomingClassGroupSessionEnrollments(r.Context(), upcoming.ID)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get upcoming class group session enrollments")
	}

	resp.EnrollmentData = append(
		make([]model.SessionEnrollment, 0, len(enrollments)),
		enrollments...,
	)
	return resp
}
