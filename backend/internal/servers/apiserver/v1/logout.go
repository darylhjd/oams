package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/oauth2"
)

type logoutResponse struct {
	response
}

func (v *APIServerV1) logout(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse = logoutResponse{newSuccessResponse()}

	oauth2.DeleteCookie(w)
	v.writeResponse(w, r, resp)
}
