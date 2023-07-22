package oauth2

import (
	"context"
	"net/url"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/darylhjd/oams/backend/internal/env"
)

const (
	mockAccessToken          = "mock-access-token"
	mockAccountHomeAccountID = "mock-home-account-id"

	MockAccountPreferredUsername = "NTU0001@e.ntu.edu.sg"
	MockIDTokenName              = "TESTACC001"
)

// MockAzureAuthenticator allows us to mock the calls to Microsoft's Azure AD APIs.
type MockAzureAuthenticator struct {
	keyCache *jwk.Cache
}

func (m *MockAzureAuthenticator) Account(_ context.Context, accountID string) (confidential.Account, error) {
	return m.mockAccount(), nil
}

func (m *MockAzureAuthenticator) RemoveAccount(_ context.Context, _ confidential.Account) error {
	return nil
}

func (m *MockAzureAuthenticator) AuthCodeURL(context.Context, string, string, []string, ...confidential.AuthCodeURLOption) (string, error) {
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

func (m *MockAzureAuthenticator) AcquireTokenByAuthCode(context.Context, string, string, []string, ...confidential.AcquireByAuthCodeOption) (confidential.AuthResult, error) {
	return m.mockAuthResult(), nil
}

func (m *MockAzureAuthenticator) AcquireTokenSilent(context.Context, []string, ...confidential.AcquireSilentOption) (confidential.AuthResult, error) {
	return m.mockAuthResult(), nil
}

func (m *MockAzureAuthenticator) GetKeyCache() *jwk.Cache {
	return m.keyCache
}

func (m *MockAzureAuthenticator) mockAuthResult() confidential.AuthResult {
	// Do this because we cannot import IDToken type.
	var result confidential.AuthResult

	result.AccessToken = mockAccessToken
	result.Account = m.mockAccount()
	result.IDToken.Name = MockIDTokenName

	return result
}

func (m *MockAzureAuthenticator) mockAccount() confidential.Account {
	return confidential.Account{
		HomeAccountID:     mockAccountHomeAccountID,
		PreferredUsername: MockAccountPreferredUsername,
	}
}

// NewMockAzureAuthenticator creates a new mock Azure Authenticator client, useful for tests.
func NewMockAzureAuthenticator() *MockAzureAuthenticator {
	return &MockAzureAuthenticator{keyCache: jwk.NewCache(context.Background())}
}
