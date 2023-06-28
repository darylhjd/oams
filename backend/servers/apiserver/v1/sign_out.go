package v1

import (
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"

	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

// signOut endpoint invalidates a user session. This is done by requesting that the browser
// remove the cookie containing the session information.
func (v *APIServerV1) signOut(w http.ResponseWriter, r *http.Request) {
	acct, ok := r.Context().Value(middleware.AccountContextKey).(confidential.Account)
	if !ok {
		http.Error(w, "unexpected account data type", http.StatusInternalServerError)
		return
	}

	if err := v.azure.RemoveAccount(r.Context(), acct); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = oauth2.DeleteSessionCookie(w)
}