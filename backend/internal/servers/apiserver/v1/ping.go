package v1

import (
	"net/http"
)

type pingResponse struct {
	response
	Message string `json:"message"`
}

// ping helps check the health of the API server. It also checks the database connection.
func (v *APIServerV1) ping(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	resp = pingResponse{
		response: newSuccessResponse(),
		Message:  "Pong~ OAMS API Service is running normally!",
	}

	if err := v.db.C.Ping(r.Context()); err != nil {
		resp = newErrorResponse(http.StatusInternalServerError, "database cannot be contacted")
	}

	v.writeResponse(w, r, resp)
}
