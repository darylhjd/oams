package v1

import (
	"net/http"
	"strconv"
	"strings"
)

func (v *APIServerV1) classGroup(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	_, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, classGroupUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, classGroupUrl, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group id"))
		return
	}

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

	v.writeResponse(w, classGroupUrl, resp)
}
