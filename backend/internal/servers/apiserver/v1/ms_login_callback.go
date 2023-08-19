package v1

import (
	"encoding/json"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

const (
	postFormCodeParam = "code"
)

// msLoginCallback handles the login callback from Microsoft Azure.
// This endpoint receives the auth code in the OAuth2 flow and exchanges this for an access token.
func (v *APIServerV1) msLoginCallback(w http.ResponseWriter, r *http.Request) {
	var s state
	if err := json.Unmarshal([]byte(r.PostFormValue(callbackStateParam)), &s); err != nil {
		v.logInternalServerError(r, err)
		v.writeResponse(w, r, newErrorResponse(http.StatusInternalServerError, "cannot parse state from login callback"))
		return
	}

	// Check that we only handle callbacks from appropriate API version.
	if s.Version != namespace {
		v.writeResponse(w, r, newErrorResponse(http.StatusTeapot, "wrong api version handling"))
		return
	}

	authResult, err := v.azure.AcquireTokenByAuthCode(
		r.Context(),
		r.PostFormValue(postFormCodeParam),
		env.GetAPIServerAzureLoginCallbackURL(),
		[]string{env.GetAPIServerAzureLoginScope()},
	)
	if err != nil {
		v.logInternalServerError(r, err)
		v.writeResponse(w, r, newErrorResponse(http.StatusInternalServerError, "cannot get auth tokens from code"))
		return
	}

	// Register user. If user has logged in before, then nothing is done.
	_, err = v.db.RegisterUser(r.Context(), database.RegisterUserParams{
		ID:    authResult.IDToken.Name,
		Email: authResult.Account.PreferredUsername,
	})
	if err != nil {
		v.logInternalServerError(r, err)
		v.writeResponse(w, r, newErrorResponse(http.StatusInternalServerError, "could not update login user details"))
		return
	}

	// Set session cookie.
	_ = oauth2.SetSessionCookie(w, authResult)

	if s.RedirectUrl == "" {
		s.RedirectUrl = Url
	}

	http.Redirect(w, r, s.RedirectUrl, http.StatusSeeOther)
}
