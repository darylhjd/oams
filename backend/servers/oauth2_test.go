package servers

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/env"
)

func Test_checkAzureToken(t *testing.T) {
	tests := []struct {
		name      string
		tokenFunc func(string) *jwt.Token
		wantErr   bool
	}{
		{
			"valid token",
			func(keyId string) *jwt.Token {
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, AzureClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    tokenIssuer,
						Audience:  jwt.ClaimStrings{env.GetAPIServerAzureClientID()},
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
						NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
					},
				})

				token.Header["kid"] = keyId
				return token
			},
			false,
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
			if err != nil {
				t.Fatal(err)
			}

			jwkKey, err := jwk.FromRaw(privateKey.PublicKey)
			if err != nil {
				t.Fatal(err)
			}

			accessToken, err := tt.tokenFunc(jwkKey.KeyID()).SignedString(privateKey)
			if err != nil {
				t.Fatal(err)
			}

			// Build request with authorization header.
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

			keySet := jwk.NewSet()
			if err = keySet.AddKey(jwkKey); err != nil {
				t.Fatal(err)
			}

			// Create handler that checks OAuth.
			_, _, err = checkAzureToken(req, keySet)
			a.Nil(err)
		})
	}
}
