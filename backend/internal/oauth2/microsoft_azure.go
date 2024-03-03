package oauth2

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
	microsoftAuthority    = fmt.Sprintf("https://login.microsoftonline.com/%s/", env.GetAPIServerAzureTenantID())
	microsoftKeySetSource = fmt.Sprintf("https://login.microsoftonline.com/%s/discovery/v2.0/keys", env.GetAPIServerAzureTenantID())
)

var (
	microsoftTokenIssuer = microsoftAuthority + "v2.0"
)

// AzureClaims is a custom struct to hold the claims received from Microsoft Azure AD.
type AzureClaims struct {
	jwt.RegisteredClaims
	AppID             string   `json:"azp"`
	Name              string   `json:"name"`
	PreferredUsername string   `json:"preferred_username"`
	Roles             []string `json:"roles"`
}

// ID returns the unique identifier for the auth user. For users, this is a normal username. For applications, this
// is a unique string identifier.
func (a AzureClaims) ID() string {
	if a.IsApplication() {
		return a.AppID
	}

	return a.Name
}

// Email returns the email for the auth user. If the user is an application, then this field will likely be empty.
func (a AzureClaims) Email() string {
	return a.PreferredUsername
}

// IsApplication checks if the auth user is an application or a normal user.
func (a AzureClaims) IsApplication() bool {
	return a.Roles != nil
}

// AppRoles returns the roles assigned to an application user. This is unused for normal users.
func (a AzureClaims) AppRoles() []string {
	return a.Roles
}

type AzureAuthenticator struct {
	Authenticator
	config *oauth2.Config
}

func (a *AzureAuthenticator) AuthCodeURL(state, verifier string) string {
	return a.config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("response_mode", "form_post"),
		oauth2.S256ChallengeOption(verifier),
	)
}

func (a *AzureAuthenticator) Exchange(ctx context.Context, code, verifier string) (*oauth2.Token, error) {
	return a.config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
}

func (a *AzureAuthenticator) CheckToken(ctx context.Context, tokenString string) (AzureClaims, *jwt.Token, error) {
	claims := AzureClaims{}

	set, err := a.GetKeyCache().Get(ctx, a.GetKeySetSource())
	if err != nil {
		return claims, nil, errors.New("could not get key set source")
	}

	token, err := parseAzureToken(set, &claims, tokenString)
	if err != nil || !token.Valid {
		return claims, nil, fmt.Errorf("token check failed: %w", err)
	}

	return claims, token, nil
}

func parseAzureToken(set jwk.Set, claims *AzureClaims, tokenString string) (*jwt.Token, error) {
	// https://learn.microsoft.com/en-us/azure/active-directory/develop/access-tokens#validating-tokens
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Header["typ"] != "JWT" {
			return nil, errors.New("wrong token type")
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not found")
		}

		keys, ok := set.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %s not found", kid)
		}

		publicKey := &rsa.PublicKey{}
		if err := keys.Raw(publicKey); err != nil {
			return nil, errors.New("could not parse key")
		}

		return publicKey, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}),
		jwt.WithAudience(env.GetAPIServerAzureClientID()),
		jwt.WithIssuer(microsoftTokenIssuer),
	)
}

// NewAzureAuthenticator creates a new Azure Authenticator for OAuth2 with Microsoft AzureAD.
func NewAzureAuthenticator() (*AzureAuthenticator, error) {
	baseAuthenticator, err := newAuthenticator(microsoftKeySetSource, microsoftTokenIssuer)
	if err != nil {
		return nil, err
	}

	return &AzureAuthenticator{*baseAuthenticator, &oauth2.Config{
		ClientID:     env.GetAPIServerAzureClientID(),
		ClientSecret: env.GetAPIServerAzureClientSecret(),
		Scopes:       []string{env.GetAPIServerAzureLoginScope()},
		Endpoint:     microsoft.AzureADEndpoint(env.GetAPIServerAzureTenantID()),
	}}, nil
}
