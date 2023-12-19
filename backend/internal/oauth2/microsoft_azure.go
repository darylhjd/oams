package oauth2

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
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
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
}

func (a AzureClaims) UserID() string {
	return a.Name
}

func (a AzureClaims) UserEmail() string {
	return a.PreferredUsername
}

type AzureAuthenticator struct {
	Authenticator
}

func (a *AzureAuthenticator) CheckToken(ctx context.Context, tokenString string) (Claims, *jwt.Token, error) {
	set, err := a.GetKeyCache().Get(ctx, a.GetKeySetSource())
	if err != nil {
		return nil, nil, errors.New("could not get key set source")
	}

	claims := AzureClaims{}

	token, err := parseAzureToken(set, &claims, tokenString)
	if err != nil || !token.Valid {
		return nil, nil, fmt.Errorf("token check failed: %w", err)
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

	return &AzureAuthenticator{Authenticator: *baseAuthenticator}, nil
}
