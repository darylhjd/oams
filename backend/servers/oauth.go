package servers

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/darylhjd/oams/backend/env"
)

const (
	AuthFieldName = "access_token"

	tokenKeySet = "https://login.microsoftonline.com/%s/discovery/v2.0/keys"
)

// checkAzureToken to make sure the token passes validation.
func checkAzureToken(r *http.Request) (*jwt.Token, error) {
	tokenHeader := r.Header.Get("Authorization")
	split := strings.Split(tokenHeader, "Bearer ")
	if len(split) < 2 {
		return nil, fmt.Errorf("malformed authorization header")
	}

	keySet, err := jwk.Fetch(r.Context(), fmt.Sprintf(tokenKeySet, env.GetAPIServerAzureTenantID()))
	if err != nil {
		return nil, fmt.Errorf("cannot fetch azure key set")
	}

	token, err := jwt.Parse(split[1], func(token *jwt.Token) (interface{}, error) {
		if token.Header["typ"] != "JWT" {
			return nil, fmt.Errorf("wrong token type")
		}

		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found")
		}

		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %s not found", kid)
		}

		publicKey := &rsa.PublicKey{}
		if err = keys.Raw(publicKey); err != nil {
			return nil, fmt.Errorf("could not parse key")
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token check failed: %w", err)
	}

	return token, nil
}
