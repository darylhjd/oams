package middleware

import (
	"context"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"

	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

const (
	AuthContextKey = "auth_context"
)

// AuthContext stores useful information regarding an authentication.
type AuthContext struct {
	Claims     *oauth2.AzureClaims
	AuthResult confidential.AuthResult
}

// GetAuthContext is a helper function to get the authentication context from a request context.
func GetAuthContext(r *http.Request) (AuthContext, bool) {
	authContext, ok := r.Context().Value(AuthContextKey).(AuthContext)
	return authContext, ok
}

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

// CheckAuthorised checks if a request is authorised for a handler and adds relevant authentication context
// if check is successful.
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

		// Add auth context to the request.
		r = r.WithContext(context.WithValue(r.Context(), AuthContextKey, AuthContext{
			Claims:     claims,
			AuthResult: res,
		}))
		handlerFunc(w, r)
	}
}
