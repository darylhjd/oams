package v1

import (
	"errors"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

type logoutResponse struct {
	response
}

func (v *APIServerV1) logout(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	authContext, isSignedIn, err := middleware.GetAuthContext(r)
	switch {
	case err == nil && !isSignedIn:
		err = errors.New("logout called but there is no session user")
		v.logInternalServerError(r, err)
		resp = newErrorResponse(http.StatusInternalServerError, err.Error())
	case err != nil:
		v.logInternalServerError(r, err)
		resp = newErrorResponse(http.StatusInternalServerError, "unexpected auth context type")
	default:
		resp = logoutResponse{newSuccessResponse()}
		if err = v.azure.RemoveAccount(r.Context(), authContext.AuthResult.Account); err != nil {
			v.logInternalServerError(r, err)
			resp = newErrorResponse(http.StatusInternalServerError, "error clearing account cache")
		}
	}

	_ = oauth2.DeleteSessionCookie(w)
	v.writeResponse(w, r, resp)
}
