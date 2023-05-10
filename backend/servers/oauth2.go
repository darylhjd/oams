package servers

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/darylhjd/oams/backend/env"
)

const (
	MicrosoftAuthority = "https://login.microsoftonline.com/%s/"
)

const (
	tokenKeySet = "https://login.microsoftonline.com/%s/discovery/v2.0/keys"
)

type AzureClaims struct {
	jwt.RegisteredClaims
	Scp   string
	Roles []string
}

// checkAzureToken to make sure the token passes validation.
func checkAzureToken(r *http.Request) (*AzureClaims, *jwt.Token, error) {
	// https://learn.microsoft.com/en-us/azure/active-directory/develop/access-tokens#validating-tokens
	tokenHeader := r.Header.Get("Authorization")
	split := strings.Split(tokenHeader, "Bearer ")
	if len(split) < 2 {
		return nil, nil, errors.New("malformed authorization header")
	}

	keySet, err := jwk.Fetch(r.Context(), fmt.Sprintf(tokenKeySet, env.GetAPIServerAzureTenantID()))
	if err != nil {
		return nil, nil, errors.New("cannot fetch azure key set")
	}

	claims := &AzureClaims{}
	issuer, err := url.JoinPath(fmt.Sprintf(MicrosoftAuthority, env.GetAPIServerAzureTenantID()), "v2.0")
	if err != nil {
		return nil, nil, errors.New("cannot build expected iss claim")
	}

	token, err := jwt.ParseWithClaims(split[1], claims, func(token *jwt.Token) (interface{}, error) {
		if token.Header["typ"] != "JWT" {
			return nil, errors.New("wrong token type")
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not found")
		}

		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %s not found", kid)
		}

		publicKey := &rsa.PublicKey{}
		if err = keys.Raw(publicKey); err != nil {
			return nil, errors.New("could not parse key")
		}

		return publicKey, nil
	},
		jwt.WithValidMethods([]string{jwa.RS256.String()}),
		jwt.WithAudience(env.GetAPIServerAzureClientID()),
		jwt.WithIssuer(issuer))
	if err != nil || !token.Valid {
		return nil, nil, fmt.Errorf("token check failed: %w", err)
	}

	return claims, token, nil
}
