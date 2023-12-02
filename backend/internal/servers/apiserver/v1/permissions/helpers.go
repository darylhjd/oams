package permissions

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/middleware/values"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

// HasPermissions checks if a user with a role has all the given permissions.
func HasPermissions(role model.UserRole, permissions ...P) bool {
	permModel, ok := rolePermissionMapping[role]
	if !ok {
		return false
	}

	for _, perm := range permissions {
		if _, ok = permModel[perm]; !ok {
			return false
		}
	}

	return true
}

// EnforceAccessPolicy based on role-based access control.
func EnforceAccessPolicy(
	handlerFunc http.HandlerFunc,
	authenticator oauth2.Authenticator,
	db *database.DB,
	methodPermissions map[string][]P,
) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		authContext := values.GetAuthContext(r.Context())

		if !HasPermissions(authContext.User.Role, methodPermissions[r.Method]...) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handlerFunc(w, r)
	}

	return middleware.MustAuth(handler, authenticator, db)
}