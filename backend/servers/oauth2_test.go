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
		name        string
		token       *jwt.Token
		wantErr     bool
		containsErr string
	}{
		{
			"valid token",
			jwt.NewWithClaims(jwt.SigningMethodRS256, AzureClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    tokenIssuer,
					Audience:  jwt.ClaimStrings{env.GetAPIServerAzureClientID()},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}),
			false,
			"",
		},
		{
			"expired token",
			jwt.NewWithClaims(jwt.SigningMethodRS256, AzureClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    tokenIssuer,
					Audience:  jwt.ClaimStrings{env.GetAPIServerAzureClientID()},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}),
			true,
			"token is expired",
		},
		{
			"token with wrong issuer",
			jwt.NewWithClaims(jwt.SigningMethodRS256, AzureClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "bad issuer",
					Audience:  jwt.ClaimStrings{env.GetAPIServerAzureClientID()},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}),
			true,
			"token has invalid issuer",
		},
		{
			"token with wrong audience",
			jwt.NewWithClaims(jwt.SigningMethodRS256, AzureClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    tokenIssuer,
					Audience:  jwt.ClaimStrings{"bad audience"},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}),
			true,
			"token has invalid audience",
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

			tt.token.Header["kid"] = jwkKey.KeyID()

			accessToken, err := tt.token.SignedString(privateKey)
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
			a.Equal(tt.wantErr, err != nil)
			if tt.wantErr {
				a.ErrorContains(err, tt.containsErr)
			}
		})
	}
}
