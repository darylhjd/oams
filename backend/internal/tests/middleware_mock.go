package tests

import (
	"github.com/darylhjd/oams/backend/internal/middleware/values"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

// NewMockAuthContext creates a middleware.AuthContext with mock information.
func NewMockAuthContext() values.AuthContext {
	return values.AuthContext{
		Claims:     &oauth2.AzureClaims{},
		AuthResult: NewMockAzureAuthenticator().MockAuthResult(),
	}
}
