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

func (v *APIServerV1) upcomingClassGroupSession(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classGroupSessionId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, upcomingClassGroupSessionUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group session id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.upcomingClassGroupSessionGet(r, classGroupSessionId)
	case http.MethodPost:
		resp = v.upcomingClassGroupSessionPost(r, classGroupSessionId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type upcomingClassGroupSessionGetResponse struct {
	response
	UpcomingClassGroupSession database.UpcomingManagedClassGroupSession `json:"upcoming_class_group_session"`
	AttendanceEntries         []database.AttendanceEntry                `json:"attendance_entries"`
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) upcomingClassGroupSessionGet(r *http.Request, sessionId int64) apiResponse {
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

	return upcomingClassGroupSessionGetResponse{
		newSuccessResponse(),
		upcoming,
		append(
			make([]database.AttendanceEntry, 0, len(entries)),
			entries...,
		),
	}
}

type upcomingClassGroupSessionPostRequest struct {
	ID            int64  `json:"id"`
	UserID        string `json:"user_id"`
	Attended      bool   `json:"attended"`
	UserSignature string `json:"user_signature"`
}

type upcomingClassGroupSessionPostResponse struct {
	response
	Attended bool `json:"attended"`
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) upcomingClassGroupSessionPost(r *http.Request, sessionId int64) apiResponse {
	var req upcomingClassGroupSessionPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	err := v.db.UpdateAttendanceEntry(r.Context(), database.UpdateAttendanceEntryParams{
		ClassGroupSessionID: sessionId,
		ID:                  req.ID,
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

	return upcomingClassGroupSessionPostResponse{
		newSuccessResponse(),
		req.Attended,
	}
}
