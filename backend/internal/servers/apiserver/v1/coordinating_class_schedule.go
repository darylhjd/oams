package v1

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/pkg/datetime"
	"github.com/darylhjd/oams/backend/pkg/to"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) coordinatingClassSchedule(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classId, err := to.Int64(r.PathValue("classId"))
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	sessionId, err := to.Int64(r.PathValue("sessionId"))
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid session id"))
		return
	}

	switch r.Method {
	case http.MethodPut:
		resp = v.coordinatingClassSchedulePut(r, classId, sessionId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type coordinatingClassSchedulePutRequest struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type coordinatingClassSchedulePutResponse struct {
	response
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (v *APIServerV1) coordinatingClassSchedulePut(r *http.Request, classId, sessionId int64) apiResponse {
	var req coordinatingClassSchedulePutRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	s, err := v.db.UpdateCoordinatingClassSchedule(r.Context(), database.UpdateCoordinatingClassScheduleParams{
		ClassID:   classId,
		SessionID: sessionId,
		StartTime: time.UnixMilli(req.StartTime).In(datetime.Location),
		EndTime:   time.UnixMilli(req.EndTime).In(datetime.Location),
	})
	if err != nil {
		switch {
		case errors.Is(err, qrm.ErrNoRows):
			return newErrorResponse(http.StatusUnauthorized, "not allowed to update coordinating class schedule")
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusBadRequest, "class group already has s at this timing")
		case database.ErrSQLState(err, database.SQLStateFailedConstraint):
			return newErrorResponse(http.StatusBadRequest, "start time must be before end time")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not update coordinating class schedule")
		}
	}

	return coordinatingClassSchedulePutResponse{
		newSuccessResponse(),
		s.StartTime,
		s.EndTime,
	}
}
