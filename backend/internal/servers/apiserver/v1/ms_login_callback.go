package v1

import (
	"encoding/json"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/oauth2"
)

const (
	callbackStateParam = "state"
	postFormCodeParam  = "code"
)

// msLoginCallback handles the login callback from Microsoft Azure.
// This endpoint receives the auth code in the OAuth2 flow and exchanges this for an access token.
func (v *APIServerV1) msLoginCallback(w http.ResponseWriter, r *http.Request) {
	var s oauthState
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

	code := r.PostFormValue(postFormCodeParam)

	token, err := v.auth.Exchange(r.Context(), code, s.Verifier)
	if err != nil {
		v.logInternalServerError(r, err)
		v.writeResponse(w, r, newErrorResponse(http.StatusInternalServerError, "could not exchange for token"))
		return
	}

	// Set session cookie.
	oauth2.SetCookie(w, token)

	if s.RedirectUrl == "" {
		s.RedirectUrl = Url
	}

	http.Redirect(w, r, s.RedirectUrl, http.StatusSeeOther)
}
