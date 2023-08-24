package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

type logoutResponse struct {
	response
}

func (v *APIServerV1) logout(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse = logoutResponse{newSuccessResponse()}

	authContext := middleware.GetAuthContext(r.Context())
	if err := v.azure.RemoveAccount(r.Context(), authContext.AuthResult.Account); err != nil {
		v.logInternalServerError(r, err)
		resp = newErrorResponse(http.StatusInternalServerError, "error clearing account cache")
	}

	_ = oauth2.DeleteSessionCookie(w)
	v.writeResponse(w, r, resp)
}
