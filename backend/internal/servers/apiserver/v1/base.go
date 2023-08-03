package v1

import (
	"net/http"
)

type baseResponse struct {
	response
	Message string `json:"message"`
}

func (v *APIServerV1) base(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	// Check if it is really the base URL.
	if r.URL.Path != baseUrl {
		resp = newErrorResponse(http.StatusNotFound, "endpoint not implemented or is not supported")
	} else {
		resp = baseResponse{
			newSuccessResponse(),
			"Welcome to OAMS API Service V1! To get started, read the API docs.",
		}
	}

	v.writeResponse(w, baseUrl, resp)
}
