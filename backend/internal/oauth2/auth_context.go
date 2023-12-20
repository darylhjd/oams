package oauth2

import (
	"context"
	"errors"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
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
	Claims Claims
	User   model.User
}

// GetAuthContext is a helper function to get the authentication context from a request context. If the auth context is
// not present, or it is of a wrong type, then GetAuthContext panics.
func GetAuthContext(ctx context.Context) AuthContext {
	val := ctx.Value(AuthContextKey)
	if val == nil {
		panic(ErrNoAuthContext)
	}

	authContext, ok := val.(AuthContext)
	if !ok {
		panic(ErrUnexpectedAuthContextType)
	}

	return authContext
}
