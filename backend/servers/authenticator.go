package servers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/darylhjd/oams/backend/env"
)

const authenticatorNamespace = "authenticator"

// Authenticator provides an interface which all authentication services for a server must implement.
type Authenticator interface {
	AuthCodeURL(ctx context.Context, clientID, redirectURI string, scopes []string, opts ...confidential.AuthCodeURLOption) (string, error)
	AcquireTokenByAuthCode(ctx context.Context, code string, redirectURI string, scopes []string, opts ...confidential.AcquireByAuthCodeOption) (confidential.AuthResult, error)
	GetKeySet() jwk.Set
}

// AzureAuthenticator is a wrapper around the Microsoft Azure AD client.
type AzureAuthenticator struct {
	*confidential.Client
	keySet jwk.Set
}

func (a *AzureAuthenticator) GetKeySet() jwk.Set {
	return a.keySet
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

	return &AzureAuthenticator{
		Client: &azureClient,
		keySet: jwk.NewCachedSet(jwk.NewCache(context.Background()), keySetSource),
	}, nil
}

// MockAzureAuthenticator allows us to mock the calls to Microsoft's Azure AD APIs.
type MockAzureAuthenticator struct {
	keySet jwk.Set
}

func (m *MockAzureAuthenticator) AuthCodeURL(ctx context.Context, clientID, redirectURI string, scopes []string, opts ...confidential.AuthCodeURLOption) (string, error) {
	path, err := url.JoinPath(env.GetAPIServerAzureTenantID(), "oauth2", "v2.0", "authorize")
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Set("client_id", env.GetAPIServerAzureClientID())
	values.Set("redirect_url", env.GetAPIServerAzureLoginCallbackURL())
	values.Set("response_type", "code")
	values.Set("scope", env.GetAPIServerAzureLoginScope())

	u := url.URL{
		Scheme:   "https",
		Host:     "login.microsoftonline.com",
		Path:     path,
		RawQuery: values.Encode(),
	}

	return u.String(), nil
}

func (m *MockAzureAuthenticator) AcquireTokenByAuthCode(
	ctx context.Context,
	code string,
	redirectURI string,
	scopes []string,
	opts ...confidential.AcquireByAuthCodeOption) (confidential.AuthResult, error) {
	return confidential.AuthResult{
		AccessToken: "mock-access-token",
	}, nil
}

func (m *MockAzureAuthenticator) GetKeySet() jwk.Set {
	return m.keySet
}

// NewMockAzureAuthenticator creates a new mock Azure Authenticator client, useful for tests.
func NewMockAzureAuthenticator() *MockAzureAuthenticator {
	return &MockAzureAuthenticator{keySet: jwk.NewSet()}
}
