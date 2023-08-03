package v1

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
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

	v.writeResponse(w, sessionEnrollmentsUrl, resp)
}

type sessionEnrollmentsGetResponse struct {
	response
	SessionEnrollments []database.SessionEnrollment `json:"session_enrollments"`
}

func (v *APIServerV1) sessionEnrollmentsGet(r *http.Request) apiResponse {
	enrollments, err := v.db.Q.ListSessionEnrollments(r.Context())
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, "could not process session enrollments get database action")
	}

	resp := sessionEnrollmentsGetResponse{
		newSuccessResponse(),
		make([]database.SessionEnrollment, 0, len(enrollments)),
	}

	resp.SessionEnrollments = append(resp.SessionEnrollments, enrollments...)
	return resp
}

type sessionEnrollmentsPostRequest struct {
	SessionEnrollment database.CreateSessionEnrollmentParams `json:"session_enrollment"`
}

type sessionEnrollmentsPostResponse struct {
	response
	SessionEnrollment database.CreateSessionEnrollmentRow `json:"session_enrollment"`
}

func (v *APIServerV1) sessionEnrollmentsPost(r *http.Request) apiResponse {
	var (
		b   bytes.Buffer
		req sessionEnrollmentsPostRequest
	)

	if _, err := b.ReadFrom(r.Body); err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if err := json.Unmarshal(b.Bytes(), &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, "could not parse request body")
	}

	session, err := v.db.Q.CreateSessionEnrollment(r.Context(), req.SessionEnrollment)
	if err != nil {
		switch {
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusConflict, "session enrollment with same session_id and user_id already exists")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusBadRequest, "session_id or user_id is not valid")
		default:
			return newErrorResponse(http.StatusInternalServerError, "could not process session enrollments post database action")
		}
	}

	return sessionEnrollmentsPostResponse{
		newSuccessResponse(),
		session,
	}
}
