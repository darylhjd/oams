package middleware

import (
	"context"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"

	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

const (
	ClaimsContextKey  = "claims"
	AccountContextKey = "account"
)

// AllowMethods allows a handler to accept only certain specified HTTP methods.
func AllowMethods(handlerFunc http.HandlerFunc, methods ...string) http.HandlerFunc {
	// Create set of allowed methods.
	allowedMethods := map[string]struct{}{}
	for _, method := range methods {
		allowedMethods[method] = struct{}{}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := allowedMethods[r.Method]; !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		handlerFunc(w, r)
	}
}

// CheckAuthorised checks if a request is authorised for a handler.
func CheckAuthorised(handlerFunc http.HandlerFunc, authenticator oauth2.Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		set, err := authenticator.GetKeyCache().Get(r.Context(), oauth2.KeySetSource)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// We will be using session cookies for authentication.
		c, err := r.Cookie(oauth2.SessionCookieIdent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// NOTE: Microsoft's authenticator library does not return an error if the account is not found in cache.
		// Instead, a zero-value Account is returned, so we check that.
		acct, err := authenticator.Account(r.Context(), c.Value)
		if err != nil || acct.IsZero() {
			http.Error(w, "account not found in session cache", http.StatusUnauthorized)
			return
		}

		// NOTE: If the backend service is restarted, all cache is lost, and all users must log in again.
		res, err := authenticator.AcquireTokenSilent(
			r.Context(),
			[]string{env.GetAPIServerAzureLoginScope()},
			confidential.WithSilentAccount(acct))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, _, err := oauth2.CheckAzureToken(set, res.AccessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Update the session cookie.
		_ = oauth2.SetSessionCookie(w, res)

		// Add relevant contexts to the request.
		r = r.WithContext(context.WithValue(r.Context(), ClaimsContextKey, claims))
		r = r.WithContext(context.WithValue(r.Context(), AccountContextKey, acct))
		handlerFunc(w, r)
	}
}
