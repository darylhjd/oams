package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

const (
	postFormCodeParam = "code"
)

// msLoginCallback handles the login callback from Microsoft Azure.
// This endpoint receives the auth code in the OAuth2 flow and exchanges this for an access token.
// TODO: Implement code challenge using PKCE.
func (v *APIServerV1) msLoginCallback(w http.ResponseWriter, r *http.Request) {
	var s state
	err := json.Unmarshal([]byte(r.PostFormValue(callbackStateParam)), &s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - cannot parse state from login callback", namespace), zap.Error(err))
		return
	}

	// Check that we only handle callbacks from appropriate API version.
	if s.Version != namespace {
		http.Error(w, "wrong API version used for login callback", http.StatusTeapot)
		v.l.Error(fmt.Sprintf("%s - received login callback of different version so ignoring", namespace),
			zap.Any(callbackStateParam, s))
		return
	}

	res, err := v.azure.AcquireTokenByAuthCode(
		r.Context(),
		r.PostFormValue(postFormCodeParam),
		env.GetAPIServerAzureLoginCallbackURL(),
		[]string{env.GetAPIServerAzureLoginScope()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not get tokens from auth code", namespace), zap.Error(err))
		return
	}

	// Set session cookie.
	_ = oauth2.SetSessionCookie(w, res)

	redirectUrl := s.ReturnTo
	if redirectUrl == "" {
		redirectUrl = Url
	}

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}
