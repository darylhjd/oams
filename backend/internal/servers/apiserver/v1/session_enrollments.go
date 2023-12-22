package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
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

	v.writeResponse(w, r, resp)
}

type sessionEnrollmentsGetResponse struct {
	response
	SessionEnrollments []model.SessionEnrollment `json:"session_enrollments"`
}

func (v *APIServerV1) sessionEnrollmentsGet(r *http.Request) apiResponse {
	params, err := database.DecodeListQueryParams(
		r.URL.Query(), table.SessionEnrollments, table.SessionEnrollments.AllColumns,
	)
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
