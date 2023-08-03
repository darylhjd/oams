package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) sessionEnrollments(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.sessionEnrollmentsGet(r)
	case http.MethodPost:
		resp = newErrorResponse(http.StatusNotImplemented, "")
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
