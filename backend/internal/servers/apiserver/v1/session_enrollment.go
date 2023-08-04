package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
		resp = v.sessionEnrollmentPatch(r, enrollmentId)
	case http.MethodDelete:
		resp = v.sessionEnrollmentDelete(r, enrollmentId)
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

type sessionEnrollmentPatchRequest struct {
	SessionEnrollment sessionEnrollmentPatchSessionEnrollmentRequestFields `json:"session_enrollment"`
}

type sessionEnrollmentPatchSessionEnrollmentRequestFields struct {
	SessionID *int64  `json:"session_id"`
	UserID    *string `json:"user_id"`
	Attended  *bool   `json:"attended"`
}

func (r sessionEnrollmentPatchRequest) updateSessionEnrollmentParams(enrollmentId int64) database.UpdateSessionEnrollmentParams {
	params := database.UpdateSessionEnrollmentParams{ID: enrollmentId}

	if r.SessionEnrollment.SessionID != nil {
		params.SessionID = pgtype.Int8{Int64: *r.SessionEnrollment.SessionID, Valid: true}
	}

	if r.SessionEnrollment.UserID != nil {
		params.UserID = pgtype.Text{String: *r.SessionEnrollment.UserID, Valid: true}
	}

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
	var (
		b   bytes.Buffer
		req sessionEnrollmentPatchRequest
	)

	if _, err := b.ReadFrom(r.Body); err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if err := json.Unmarshal(b.Bytes(), &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, "could not parse request body")
	}

	enrollment, err := v.db.Q.UpdateSessionEnrollment(r.Context(), req.updateSessionEnrollmentParams(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "session enrollment to update does not exist")
		}

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
	_, err := v.db.Q.DeleteSessionEnrollment(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "session enrollment to delete does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process session enrollment delete database action")
	}

	return sessionEnrollmentDeleteResponse{newSuccessResponse()}
}
