package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/jackc/pgx/v5"
)

func (v *APIServerV1) sessionEnrollment(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	enrollmentId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, sessionEnrollmentUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, sessionEnrollmentUrl, newErrorResponse(http.StatusUnprocessableEntity, "invalid session_enrollment id"))
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.sessionEnrollmentGet(r, enrollmentId)
	case http.MethodPatch:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, sessionEnrollmentUrl, resp)
}

type sessionEnrollmentGetResponse struct {
	response
	SessionEnrollment database.SessionEnrollment `json:"session_enrollment"`
}

func (v *APIServerV1) sessionEnrollmentGet(r *http.Request, id int64) apiResponse {
	enrollment, err := v.db.Q.GetSessionEnrollment(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested session enrollment does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process session enrollment get database action")
	}

	return sessionEnrollmentGetResponse{
		newSuccessResponse(),
		enrollment,
	}
}
