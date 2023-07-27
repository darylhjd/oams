package v1

import (
	"net/http"
	"strings"
)

func (v *APIServerV1) class(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	_ = strings.TrimPrefix(r.URL.Path, classUrl)
	switch r.Method {
	case http.MethodGet:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodPut:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")

	}

	v.writeResponse(w, classUrl, resp)
}
