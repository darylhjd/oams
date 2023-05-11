package servers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
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

func TestAllowMethods(t *testing.T) {
	tests := []struct {
		name           string
		testMethods    []string
		allowedMethods []string
		wantErr        bool
	}{
		{
			"success on allowed methods",
			[]string{
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
				http.MethodConnect,
				http.MethodHead,
				http.MethodOptions,
				http.MethodPatch,
				http.MethodTrace,
				http.MethodPut},
			[]string{
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
				http.MethodConnect,
				http.MethodHead,
				http.MethodOptions,
				http.MethodPatch,
				http.MethodTrace,
				http.MethodPut,
			},
			false,
		},
		{
			"fail on disallowed method",
			[]string{http.MethodPost},
			[]string{http.MethodGet},
			true,
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, method := range tt.testMethods {
				req, err := http.NewRequest(method, "", nil)
				if err != nil {
					t.Fatal(err)
				}

				testHandler := AllowMethods(func(w http.ResponseWriter, r *http.Request) {}, tt.allowedMethods...)
				rr := httptest.NewRecorder()

				testHandler.ServeHTTP(rr, req)

				a.Equal(tt.wantErr, rr.Result().StatusCode == http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestCheckAuthorised(t *testing.T) {
	tests := []struct {
		name      string
		tokenFunc func(*rsa.PublicKey) *jwt.Token
		wantErr   bool
	}{
		{
			"valid token",
			func(pubKey *rsa.PublicKey) *jwt.Token {
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, AzureClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    tokenIssuer,
						Audience:  jwt.ClaimStrings{env.GetAPIServerAzureClientID()},
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
						NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
					},
				})

				token.Header["kid"] = x509.MarshalPKCS1PublicKey(pubKey)
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

			accessToken, err := tt.tokenFunc(&privateKey.PublicKey).SignedString(privateKey)
			if err != nil {
				t.Fatal(err)
			}

			// Build request with authorization header.
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

			keySet := jwk.NewSet()
			jwkKey, err := jwk.FromRaw(privateKey.PublicKey)
			if err != nil {
				t.Fatal(err)
			}

			if err = keySet.AddKey(jwkKey); err != nil {
				t.Fatal(err)
			}

			// Create handler that checks OAuth.
			handler := CheckAuthorised(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}, keySet)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			t.Log(rr.Result().StatusCode)
			a.Equal(tt.wantErr, rr.Result().StatusCode != http.StatusOK)
		})
	}
}
