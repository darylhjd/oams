package middleware

import (
	"context"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/middleware/values"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	"github.com/darylhjd/oams/backend/internal/permissions"
)

// AllowMethods allows a handler to accept only certain specified HTTP methods.
func AllowMethods(handlerFunc http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, method := range methods {
			if method == r.Method {
				handlerFunc(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// AllowMethodsWithPermissions allows a handler to accept only requests from users with certain permissions.
func AllowMethodsWithPermissions(handlerFunc http.HandlerFunc, db *database.DB, methodPermissions map[string][]permissions.Permission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authContext := values.GetAuthContext(r.Context())

		for method, roles := range methodPermissions {
			if method == r.Method {
				user, err := db.GetUser(r.Context(), authContext.AuthResult.IDToken.Name)
				if err != nil {
					http.Error(w, "error getting auth user", http.StatusInternalServerError)
					return
				}

				if !permissions.HasPermissions(user.Role, roles...) {
					w.WriteHeader(http.StatusUnauthorized)
					return

				}

				handlerFunc(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// MustAuth adds AuthContext for a handler and checks for authentication status.
func MustAuth(handlerFunc http.HandlerFunc, authenticator oauth2.Authenticator) http.HandlerFunc {
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
		r = r.WithContext(context.WithValue(r.Context(), values.AuthContextKey, values.AuthContext{
			Claims:     claims,
			AuthResult: res,
		}))
		handlerFunc(w, r)
	}
}
