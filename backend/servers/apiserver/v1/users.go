package v1

import (
	"fmt"
	"net/http"
)

// users endpoint returns useful information on the current session user and information on any requested users.
func (v *APIServerV1) users(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.usersGet(r)
	case http.MethodPut:
		resp = v.usersPut(r)
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, fmt.Sprintf("method %s is not allowed", r.Method))
	}

	v.writeResponse(w, usersUrl, resp)
}
