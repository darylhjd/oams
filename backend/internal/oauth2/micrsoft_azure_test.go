package oauth2

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/env"
)

func Test_checkAzureToken(t *testing.T) {
	tests := []struct {
		name    string
		token   *jwt.Token
		wantErr error
	}{
		{
			"valid token",
			jwt.NewWithClaims(jwt.SigningMethodRS256, AzureClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    microsoftTokenIssuer,
					Audience:  jwt.ClaimStrings{env.GetAPIServerAzureClientID()},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}),
			nil,
		},
		{
			"expired token",
			jwt.NewWithClaims(jwt.SigningMethodRS256, AzureClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    microsoftTokenIssuer,
					Audience:  jwt.ClaimStrings{env.GetAPIServerAzureClientID()},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}),
			jwt.ErrTokenExpired,
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
			jwt.ErrTokenInvalidIssuer,
		},
		{
			"token with wrong audience",
			jwt.NewWithClaims(jwt.SigningMethodRS256, AzureClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    microsoftTokenIssuer,
					Audience:  jwt.ClaimStrings{"bad audience"},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}),
			jwt.ErrTokenInvalidAudience,
		},
		{
			"token with wrong signing method",
			jwt.NewWithClaims(jwt.SigningMethodRS384, AzureClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    microsoftTokenIssuer,
					Audience:  jwt.ClaimStrings{env.GetAPIServerAzureClientID()},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}),
			jwt.ErrTokenSignatureInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

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

			keySet := jwk.NewSet()
			if err = keySet.AddKey(jwkKey); err != nil {
				t.Fatal(err)
			}

			claims := &AzureClaims{}
			_, err = parseAzureToken(keySet, claims, accessToken)
			a.ErrorIs(err, tt.wantErr)
		})
	}
}
