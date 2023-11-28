package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
)

func (v *APIServerV1) sessionEnrollments(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.sessionEnrollmentsGet(r)
	case http.MethodPost:
		resp = v.sessionEnrollmentsPost(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type sessionEnrollmentsGetResponse struct {
	response
	SessionEnrollments []model.SessionEnrollment `json:"session_enrollments"`
}

func (v *APIServerV1) sessionEnrollmentsGet(r *http.Request) apiResponse {
	params, err := database.DecodeListQueryParams(r.URL.Query(), table.SessionEnrollments.AllColumns)
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}

	enrollments, err := v.db.ListSessionEnrollments(r.Context(), params)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process session enrollments get database action")
	}

	resp := sessionEnrollmentsGetResponse{
		newSuccessResponse(),
		make([]model.SessionEnrollment, 0, len(enrollments)),
	}

	resp.SessionEnrollments = append(resp.SessionEnrollments, enrollments...)
	return resp
}

type sessionEnrollmentsPostRequest struct {
	SessionEnrollment database.CreateSessionEnrollmentParams `json:"session_enrollment"`
}

type sessionEnrollmentsPostResponse struct {
	response
	SessionEnrollment sessionEnrollmentsPostSessionEnrollmentResponseFields `json:"session_enrollment"`
}

type sessionEnrollmentsPostSessionEnrollmentResponseFields struct {
	ID        int64     `json:"id"`
	SessionID int64     `json:"session_id"`
	UserID    string    `json:"user_id"`
	Attended  bool      `json:"attended"`
	CreatedAt time.Time `json:"created_at"`
}

func (v *APIServerV1) sessionEnrollmentsPost(r *http.Request) apiResponse {
	var req sessionEnrollmentsPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	session, err := v.db.CreateSessionEnrollment(r.Context(), req.SessionEnrollment)
	if err != nil {
		switch {
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusConflict, "session enrollment with same session_id and user_id already exists")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusBadRequest, "session_id and/or user_id does not exist")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process session enrollments post database action")
		}
	}

	return sessionEnrollmentsPostResponse{
		newSuccessResponse(),
		sessionEnrollmentsPostSessionEnrollmentResponseFields{
			session.ID,
			session.SessionID,
			session.UserID,
			session.Attended,
			session.CreatedAt,
		},
	}
}
