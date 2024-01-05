package v1

import (
	"errors"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) upcomingClassGroupSessionAttendances(w http.ResponseWriter, r *http.Request, sessionId int64) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.upcomingClassGroupSessionAttendancesGet(r, sessionId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type upcomingClassGroupSessionAttendancesGetResponse struct {
	response
	UpcomingClassGroupSession database.UpcomingManagedClassGroupSession `json:"upcoming_class_group_session"`
	AttendanceEntries         []database.AttendanceEntry                `json:"attendance_entries"`
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) upcomingClassGroupSessionAttendancesGet(r *http.Request, sessionId int64) apiResponse {
	upcoming, err := v.db.GetUpcomingManagedClassGroupSession(r.Context(), sessionId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested upcoming class group session does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process upcoming class group session get database action")
	}

	entries, err := v.db.GetUpcomingClassGroupAttendanceEntries(r.Context(), upcoming.ID)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get upcoming class group session attendance entries")
	}

	return upcomingClassGroupSessionAttendancesGetResponse{
		newSuccessResponse(),
		upcoming,
		append(
			make([]database.AttendanceEntry, 0, len(entries)),
			entries...,
		),
	}
}
