package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	"github.com/darylhjd/oams/backend/internal/permissions"
)

const (
	AuthContextKey = "auth_context"
)

var (
	ErrNoAuthContext             = errors.New("no auth context found")
	ErrUnexpectedAuthContextType = errors.New("unexpected auth context type")
)

// AuthContext stores useful information regarding an authentication.
type AuthContext struct {
	Claims     *oauth2.AzureClaims
	AuthResult confidential.AuthResult
}

// GetAuthContext is a helper function to get the authentication context from a request context.
func GetAuthContext(r *http.Request) (AuthContext, error) {
	val := r.Context().Value(AuthContextKey)
	if val == nil {
		return AuthContext{}, ErrNoAuthContext
	}

	authContext, ok := val.(AuthContext)
	if !ok {
		return AuthContext{}, ErrUnexpectedAuthContextType
	}

	return authContext, nil
}

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

// AllowMethodsWithUserRoles allows a handler to accept only certain specified HTTP methods with corresponding user roles.
func AllowMethodsWithUserRoles(handlerFunc http.HandlerFunc, db *database.DB, methodPermissions map[string][]permissions.Permission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authContext, err := GetAuthContext(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for method, roles := range methodPermissions {
			if method == r.Method {
				user, err := db.GetUser(r.Context(), authContext.AuthResult.IDToken.Name)
				if err != nil {
					http.Error(w, "error getting auth user", http.StatusInternalServerError)
					return
				}

				if !permissions.HasPermission(user.Role, roles...) {
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
		r = r.WithContext(context.WithValue(r.Context(), AuthContextKey, AuthContext{
			Claims:     claims,
			AuthResult: res,
		}))
		handlerFunc(w, r)
	}
}
