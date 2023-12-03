package tests

import (
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/middleware/values"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

// StubAuthContext creates a middleware.AuthContext with mock information.
// Note that the user associated with this auth context is not created in the database.
func StubAuthContext() values.AuthContext {
	return values.AuthContext{
		Claims:     &oauth2.AzureClaims{},
		AuthResult: NewMockAzureAuthenticator().MockAuthResult(),
		User: model.User{
			ID:    MockAuthenticatorUserID,
			Email: MockAuthenticatorUserEmail,
			Role:  MockAuthenticatorUserRole,
		},
	}
}
