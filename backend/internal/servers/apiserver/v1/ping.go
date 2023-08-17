package v1

import (
	"net/http"
)

type pingResponse struct {
	response
	Message string `json:"message"`
}

func (v *APIServerV1) ping(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	resp = pingResponse{
		response: newSuccessResponse(),
		Message:  "Pong~ OAMS API Service is running normally!",
	}

	if err := v.db.Conn.PingContext(r.Context()); err != nil {
		v.logInternalServerError(r, err)
		resp = newErrorResponse(http.StatusInternalServerError, "database cannot be contacted")
	}

	v.writeResponse(w, r, resp)
}
