package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) upcomingClassGroupSessionAttendance(w http.ResponseWriter, r *http.Request, sessionId int64) {
	var resp apiResponse

	enrollmentId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, upcomingClassGroupSessionAttendanceUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid session enrollment id"))
		return
	}

	switch r.Method {
	case http.MethodPatch:
		resp = v.upcomingClassGroupSessionAttendancePatch(r, sessionId, enrollmentId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type upcomingClassGroupSessionAttendancePatchRequest struct {
	Attended      bool   `json:"attended"`
	UserID        string `json:"user_id"`
	UserSignature string `json:"user_signature"`
}

type upcomingClassGroupSessionAttendancePatchResponse struct {
	response
	Attended bool `json:"attended"`
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) upcomingClassGroupSessionAttendancePatch(r *http.Request, sessionId, enrollmentId int64) apiResponse {
	var req upcomingClassGroupSessionAttendancePatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	err := v.db.UpdateAttendanceEntry(r.Context(), database.UpdateAttendanceEntryParams{
		ClassGroupSessionID: sessionId,
		SessionEnrollmentID: enrollmentId,
		UserID:              req.UserID,
		Attended:            req.Attended,
		UserSignature:       req.UserSignature,
	})
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusBadRequest, "not allowed to take attendance")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not update attendance")
	}

	return upcomingClassGroupSessionAttendancePatchResponse{
		newSuccessResponse(),
		req.Attended,
	}
}