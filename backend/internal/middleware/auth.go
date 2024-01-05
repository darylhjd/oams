package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	"github.com/golang-jwt/jwt/v5"
)

// MustAuth adds AuthContext for a handler and checks for authentication status.
func MustAuth(handlerFunc http.HandlerFunc, auth oauth2.AuthProvider, db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, _, err := CheckAuthorizationToken(r, auth)
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

// CheckAuthorizationToken to see if request is paired with a valid user session.
func CheckAuthorizationToken(r *http.Request, auth oauth2.AuthProvider) (oauth2.Claims, *jwt.Token, error) {
	accessToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	return auth.CheckToken(r.Context(), accessToken)
}
