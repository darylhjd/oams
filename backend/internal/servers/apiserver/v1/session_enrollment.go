package v1

import (
	"errors"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/pkg/to"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) sessionEnrollment(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	enrollmentId, err := to.Int64(r.PathValue("enrollmentId"))
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid session enrollment id"))
		return
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

	v.writeResponse(w, r, resp)
}

type sessionEnrollmentGetResponse struct {
	response
	SessionEnrollment model.SessionEnrollment `json:"session_enrollment"`
}

func (v *APIServerV1) sessionEnrollmentGet(r *http.Request, enrollmentId int64) apiResponse {
	enrollment, err := v.db.GetSessionEnrollment(r.Context(), enrollmentId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested session enrollment does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process session enrollment get database action")
	}

	return sessionEnrollmentGetResponse{
		newSuccessResponse(),
		enrollment,
	}
}
