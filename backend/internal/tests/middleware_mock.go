package tests

import (
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

// NewMockAuthContext creates a middleware.AuthContext with mock information.
func NewMockAuthContext() middleware.AuthContext {
	return middleware.AuthContext{
		Claims:     &oauth2.AzureClaims{},
		AuthResult: NewMockAzureAuthenticator().MockAuthResult(),
	}
}
