package v1

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

const namespace = "apiserver/v1"

const (
	Url = "/api/v1/"
)

const (
	baseUrl            = "/"
	pingUrl            = "/ping"
	loginUrl           = "/login"
	msLoginCallbackUrl = "/ms-login-callback"
	protectedUrl       = "/protected"
)

type APIServerV1 struct {
	l   *zap.Logger
	db  *database.DB
	mux *http.ServeMux

	azure oauth2.Authenticator
}

// NewAPIServerV1 creates a new APIServerV1. This is a sub-router and should not be used as a base router.
func NewAPIServerV1(l *zap.Logger, db *database.DB, azureClient oauth2.Authenticator) *APIServerV1 {
	server := APIServerV1{l, db, http.NewServeMux(), azureClient}
	server.registerHandlers()

	return &server
}

func (v *APIServerV1) registerHandlers() {
	v.mux.HandleFunc(baseUrl, middleware.AllowMethods(v.base, http.MethodGet))
	v.mux.HandleFunc(pingUrl, middleware.AllowMethods(v.ping, http.MethodGet))
	v.mux.HandleFunc(loginUrl, v.login)
	v.mux.HandleFunc(msLoginCallbackUrl, middleware.AllowMethods(v.msLoginCallback, http.MethodPost))
	v.mux.HandleFunc(protectedUrl, middleware.CheckAuthorised(v.protected, v.azure))
}

func (v *APIServerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.mux.ServeHTTP(w, r)
}
