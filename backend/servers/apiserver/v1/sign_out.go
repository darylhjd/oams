package v1

import (
	"fmt"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

// signOut endpoint invalidates a users session. This is done by requesting that the browser
// remove the cookie containing the session information.
func (v *APIServerV1) signOut(w http.ResponseWriter, r *http.Request) {
	authContext, isSignedIn, err := middleware.GetAuthContext(r)
	switch {
	case err != nil:
		break
	case !isSignedIn:
		err = fmt.Errorf("%s - sign out endpoint called with no user session", namespace)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := v.azure.RemoveAccount(r.Context(), authContext.AuthResult.Account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = oauth2.DeleteSessionCookie(w)
}
