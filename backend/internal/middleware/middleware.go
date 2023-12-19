package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

// MustAuth adds AuthContext for a handler and checks for authentication status.
func MustAuth(handlerFunc http.HandlerFunc, auth oauth2.AuthProvider, db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		claims, _, err := auth.CheckToken(r.Context(), accessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := db.RegisterUser(r.Context(), database.RegisterUserParams{
			ID:    claims.UserID(),
			Email: claims.UserEmail(),
		})
		if err != nil {
			http.Error(w, "could not get user information", http.StatusInternalServerError)
			return
		}

		// Add auth context to the request.
		r = r.WithContext(context.WithValue(r.Context(), oauth2.AuthContextKey, oauth2.AuthContext{
			Claims: claims,
			User:   user,
		}))
		handlerFunc(w, r)
	}
}
