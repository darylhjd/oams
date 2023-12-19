package tests

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/oauth2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type MockAuthenticator struct {
	keyCache *jwk.Cache
}

func (m *MockAuthenticator) GetKeyCache() *jwk.Cache {
	return m.keyCache
}

func (m *MockAuthenticator) GetKeySetSource() string {
	return uuid.NewString()
}

func (m *MockAuthenticator) GetIssuer() string {
	return uuid.NewString()
}

func (m *MockAuthenticator) CheckToken(_ context.Context, _ string) (oauth2.Claims, *jwt.Token, error) {
	return nil, nil, nil
}

// NewMockAzureAuthenticator creates a new mock Azure Authenticator client, useful for tests.
func NewMockAzureAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{keyCache: jwk.NewCache(context.Background())}
}
