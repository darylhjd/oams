package servers

import (
	"context"
	"net/url"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"

	"github.com/darylhjd/oams/backend/env"
)

// Authenticator provides an interface which all authentication services for a server must implement.
type Authenticator interface {
	AuthCodeURL(ctx context.Context, clientID, redirectURI string, scopes []string, opts ...confidential.AuthCodeURLOption) (string, error)
	AcquireTokenByAuthCode(ctx context.Context, code string, redirectURI string, scopes []string, opts ...confidential.AcquireByAuthCodeOption) (confidential.AuthResult, error)
}

// MockAzureClient allows us to mock the calls to Microsoft's APIs.
type MockAzureClient struct{}

func (m *MockAzureClient) AuthCodeURL(ctx context.Context, clientID, redirectURI string, scopes []string, opts ...confidential.AuthCodeURLOption) (string, error) {
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

func (m *MockAzureClient) AcquireTokenByAuthCode(
	ctx context.Context,
	code string,
	redirectURI string,
	scopes []string,
	opts ...confidential.AcquireByAuthCodeOption) (confidential.AuthResult, error) {
	return confidential.AuthResult{
		AccessToken: "mock-access-token",
	}, nil
}
