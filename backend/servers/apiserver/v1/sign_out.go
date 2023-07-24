package v1

import (
	"errors"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

type signOutResponse struct {
	response
}

// signOut endpoint invalidates a user session. This is done by requesting that the browser
// remove the cookie containing the session information.
func (v *APIServerV1) signOut(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	authContext, isSignedIn, err := middleware.GetAuthContext(r)
	switch {
	case err == nil && !isSignedIn:
		err = errors.New("sign-out called but there is no session user")
		fallthrough
	case err != nil:
		resp = newErrorResponse(http.StatusInternalServerError, err.Error())
	default:
		resp = signOutResponse{newSuccessfulResponse()}
		if err = v.azure.RemoveAccount(r.Context(), authContext.AuthResult.Account); err != nil {
			resp = newErrorResponse(http.StatusInternalServerError, err.Error())
		}
	}

	_ = oauth2.DeleteSessionCookie(w)
	v.writeResponse(w, signOutUrl, resp)
}
