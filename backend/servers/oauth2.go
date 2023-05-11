package servers

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/darylhjd/oams/backend/env"
)

var (
	MicrosoftAuthority = fmt.Sprintf("https://login.microsoftonline.com/%s/", env.GetAPIServerAzureTenantID())
)

var (
	tokenIssuer  = MicrosoftAuthority + "v2.0"
	keySetSource = fmt.Sprintf("https://login.microsoftonline.com/%s/discovery/v2.0/keys", env.GetAPIServerAzureTenantID())
)

// AzureClaims is a custom struct to hold the claims received from Microsoft Azure AD.
type AzureClaims struct {
	jwt.RegisteredClaims
	Scp   string
	Roles []string
}

// checkAzureToken to make sure the token passes validation.
func checkAzureToken(r *http.Request, set jwk.Set) (*AzureClaims, *jwt.Token, error) {
	// https://learn.microsoft.com/en-us/azure/active-directory/develop/access-tokens#validating-tokens
	tokenHeader := r.Header.Get("Authorization")
	split := strings.Split(tokenHeader, "Bearer ")
	if len(split) < 2 {
		return nil, nil, errors.New("malformed authorization header")
	}

	claims := &AzureClaims{}
	token, err := jwt.ParseWithClaims(split[1], claims, func(token *jwt.Token) (interface{}, error) {
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
		jwt.WithIssuer(tokenIssuer))
	if err != nil || !token.Valid {
		return nil, nil, fmt.Errorf("token check failed: %w", err)
	}

	return claims, token, nil
}
