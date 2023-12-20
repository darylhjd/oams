package tests

import (
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

const (
	MockAuthenticatorUserID    = "TESTACC001"
	MockAuthenticatorUserEmail = "NTU0001@e.ntu.edu.sg"
	MockAuthenticatorUserRole  = model.UserRole_SystemAdmin
)

// StubAuthContext creates a middleware.AuthContext with mock information.
// Note that the user associated with this auth context is not created in the database.
func StubAuthContext() oauth2.AuthContext {
	return oauth2.AuthContext{
		Claims: oauth2.AzureClaims{},
		User: model.User{
			ID:    MockAuthenticatorUserID,
			Email: MockAuthenticatorUserEmail,
			Role:  MockAuthenticatorUserRole,
		},
	}
}
