package oauth2

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

const authenticatorNamespace = "authenticator"

type AuthProvider interface {
	GetKeyCache() *jwk.Cache
	GetKeySetSource() string
	GetIssuer() string
	CheckToken(context.Context, string) (Claims, *jwt.Token, error)
}

// Authenticator provides a key cache that is used for verifying access tokens.
type Authenticator struct {
	keyCache     *jwk.Cache
	keySetSource string
	issuer       string
}

func (a *Authenticator) GetKeyCache() *jwk.Cache {
	return a.keyCache
}

func (a *Authenticator) GetKeySetSource() string {
	return a.keySetSource
}

func (a *Authenticator) GetIssuer() string {
	return a.issuer
}

func newAuthenticator(keySetSource string, issuer string) (*Authenticator, error) {
	cache := jwk.NewCache(context.Background())
	if err := cache.Register(keySetSource); err != nil {
		return nil, fmt.Errorf("%s - could not register jwt key set source: %w", authenticatorNamespace, err)
	}

	// Refresh once to fail early.
	if _, err := cache.Refresh(context.Background(), keySetSource); err != nil {
		return nil, fmt.Errorf("%s - cache source not reachable: %w", authenticatorNamespace, err)
	}

	return &Authenticator{
		keyCache:     cache,
		keySetSource: keySetSource,
		issuer:       issuer,
	}, nil
}
