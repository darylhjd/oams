package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (v *APIServerV1) sessionEnrollment(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	enrollmentId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, sessionEnrollmentUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid session enrollment id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.sessionEnrollmentGet(r, enrollmentId)
	case http.MethodPatch:
		resp = v.sessionEnrollmentPatch(r, enrollmentId)
	case http.MethodDelete:
		resp = v.sessionEnrollmentDelete(r, enrollmentId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type sessionEnrollmentGetResponse struct {
	response
	SessionEnrollment model.SessionEnrollment `json:"session_enrollment"`
}

func (v *APIServerV1) sessionEnrollmentGet(r *http.Request, id int64) apiResponse {
	enrollment, err := v.db.GetSessionEnrollment(r.Context(), id)
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

type sessionEnrollmentPatchRequest struct {
	SessionEnrollment sessionEnrollmentPatchSessionEnrollmentRequestFields `json:"session_enrollment"`
}

type sessionEnrollmentPatchSessionEnrollmentRequestFields struct {
	Attended *bool `json:"attended"`
}

func (r sessionEnrollmentPatchRequest) updateSessionEnrollmentParams(enrollmentId int64) database.UpdateSessionEnrollmentParams {
	params := database.UpdateSessionEnrollmentParams{ID: enrollmentId}

	if r.SessionEnrollment.Attended != nil {
		params.Attended = pgtype.Bool{Bool: *r.SessionEnrollment.Attended, Valid: true}
	}

	return params
}

type sessionEnrollmentPatchResponse struct {
	response
	SessionEnrollment database.UpdateSessionEnrollmentRow `json:"session_enrollment"`
}

func (v *APIServerV1) sessionEnrollmentPatch(r *http.Request, id int64) apiResponse {
	var req sessionEnrollmentPatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	enrollment, err := v.db.Q.UpdateSessionEnrollment(r.Context(), req.updateSessionEnrollmentParams(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "session enrollment to update does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process session enrollment patch database action")
	}

	return sessionEnrollmentPatchResponse{
		newSuccessResponse(),
		enrollment,
	}
}

type sessionEnrollmentDeleteResponse struct {
	response
}

func (v *APIServerV1) sessionEnrollmentDelete(r *http.Request, id int64) apiResponse {
	_, err := v.db.DeleteSessionEnrollment(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "session enrollment to delete does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process session enrollment delete database action")
	}

	return sessionEnrollmentDeleteResponse{newSuccessResponse()}
}
