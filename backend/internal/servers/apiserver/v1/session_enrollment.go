package v1

import (
	"net/http"
	"strconv"
	"strings"
)

func (v *APIServerV1) sessionEnrollment(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	_, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, sessionEnrollmentUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, sessionEnrollmentUrl, newErrorResponse(http.StatusUnprocessableEntity, "invalid session_enrollment id"))
	}

	switch r.Method {
	case http.MethodGet:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodPatch:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, sessionEnrollmentUrl, resp)
}
