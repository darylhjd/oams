package values

import (
	"context"
	"errors"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/oauth2"
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
	User       model.User
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
