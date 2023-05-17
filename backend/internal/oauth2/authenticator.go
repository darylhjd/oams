package oauth2

import (
	"context"
	"fmt"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/darylhjd/oams/backend/internal/env"
)

const authenticatorNamespace = "authenticator"

// Authenticator provides an interface which all authentication services for a server must implement.
type Authenticator interface {
	Account(ctx context.Context, accountID string) (confidential.Account, error)
	AuthCodeURL(ctx context.Context, clientID, redirectURI string, scopes []string, opts ...confidential.AuthCodeURLOption) (string, error)
	AcquireTokenByAuthCode(ctx context.Context, code string, redirectURI string, scopes []string, opts ...confidential.AcquireByAuthCodeOption) (confidential.AuthResult, error)
	AcquireTokenSilent(ctx context.Context, scopes []string, opts ...confidential.AcquireSilentOption) (confidential.AuthResult, error)
	GetKeyCache() *jwk.Cache
}

// AzureAuthenticator is a wrapper around the Microsoft Azure AD client.
type AzureAuthenticator struct {
	*confidential.Client
	keyCache *jwk.Cache
}

func (a *AzureAuthenticator) GetKeyCache() *jwk.Cache {
	return a.keyCache
}

// NewAzureAuthenticator creates a new Azure Authenticator.
func NewAzureAuthenticator() (*AzureAuthenticator, error) {
	cred, err := confidential.NewCredFromSecret(env.GetAPIServerAzureClientSecret())
	if err != nil {
		return nil, fmt.Errorf("%s - could not create credential from client secret: %w", authenticatorNamespace, err)
	}

	azureClient, err := confidential.New(MicrosoftAuthority, env.GetAPIServerAzureClientID(), cred)
	if err != nil {
		return nil, fmt.Errorf("%s - could not create azure client: %w", authenticatorNamespace, err)
	}

	cache := jwk.NewCache(context.Background())
	if err = cache.Register(KeySetSource); err != nil {
		return nil, fmt.Errorf("%s - could not register jwt key set source: %w", authenticatorNamespace, err)
	}

	// Refresh once to fail early.
	if _, err = cache.Refresh(context.Background(), KeySetSource); err != nil {
		return nil, fmt.Errorf("%s - cache source not reachable: %w", authenticatorNamespace, err)
	}

	return &AzureAuthenticator{
		Client:   &azureClient,
		keyCache: cache,
	}, nil
}
