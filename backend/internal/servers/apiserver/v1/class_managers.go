package v1

import (
	"net/http"
)

func (v *APIServerV1) classManagers(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodPost:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}
