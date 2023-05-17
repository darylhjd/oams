package v1

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

const (
	namespace = "apiserver/v1"
	Url       = "/api/v1/"
)

const (
	baseUrl            = "/"
	pingUrl            = "/ping"
	loginUrl           = "/login"
	msLoginCallbackUrl = "/ms-login-callback"
	protectedUrl       = "/protected"
)

type APIServerV1 struct {
	l  *zap.Logger
	db *database.DB

	azure oauth2.Authenticator
}

func (v *APIServerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc(baseUrl, middleware.AllowMethods(v.base, http.MethodGet))
	mux.HandleFunc(pingUrl, middleware.AllowMethods(v.ping, http.MethodGet))
	mux.HandleFunc(loginUrl, v.login)
	mux.HandleFunc(msLoginCallbackUrl, middleware.AllowMethods(v.msLoginCallback, http.MethodPost))
	mux.HandleFunc(protectedUrl, middleware.CheckAuthorised(v.protected, v.azure))

	mux.ServeHTTP(w, r)
}

// NewAPIServerV1 creates a new APIServerV1. This is a sub-router and should not be used
// as a base router.
func NewAPIServerV1(l *zap.Logger, db *database.DB, azureClient oauth2.Authenticator) *APIServerV1 {
	return &APIServerV1{l, db, azureClient}
}
