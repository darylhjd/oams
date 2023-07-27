package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"go.uber.org/zap"

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
		v.writeResponse(w, msLoginCallbackUrl, newErrorResponse(http.StatusInternalServerError, "cannot parse state from login callback"))
		return
	}

	// Check that we only handle callbacks from appropriate API version.
	if s.Version != namespace {
		v.writeResponse(w, msLoginCallbackUrl, newErrorResponse(http.StatusTeapot, "wrong api version handling"))
		return
	}

	authResult, err := v.azure.AcquireTokenByAuthCode(
		r.Context(),
		r.PostFormValue(postFormCodeParam),
		env.GetAPIServerAzureLoginCallbackURL(),
		[]string{env.GetAPIServerAzureLoginScope()},
	)
	if err != nil {
		v.writeResponse(w, msLoginCallbackUrl, newErrorResponse(http.StatusInternalServerError, "cannot get auth tokens from code"))
		return
	}

	// Upsert user into database.
	err = v.db.Q.UpsertUsers(r.Context(), []database.UpsertUsersParams{
		{
			ID:    authResult.IDToken.Name,
			Email: authResult.Account.PreferredUsername,
			// TODO: Set correct role based on auth result.
			Role: database.UserRoleSTUDENT,
		},
	}).Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - error upserting users info", namespace), zap.Error(err))
		return
	}

	// Set session cookie.
	_ = oauth2.SetSessionCookie(w, authResult)

	if s.RedirectUrl == "" {
		s.RedirectUrl = Url
	}

	http.Redirect(w, r, s.RedirectUrl, http.StatusSeeOther)
}
