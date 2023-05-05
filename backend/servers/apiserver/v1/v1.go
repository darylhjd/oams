package v1

import (
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/database"
	"github.com/darylhjd/oams/backend/servers"
)

const (
	namespace = "apiserver/v1"

	baseUrl            = "/"
	pingUrl            = "/ping"
	loginUrl           = "/login"
	msLoginCallbackUrl = "/ms-login-callback"
)

type APIServerV1 struct {
	l  *zap.Logger
	db *database.DB

	azure *confidential.Client
}

func (v *APIServerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc(baseUrl, servers.AllowMethods(v.base, http.MethodGet))
	mux.HandleFunc(pingUrl, servers.AllowMethods(v.ping, http.MethodGet))
	mux.HandleFunc(loginUrl, v.login)
	mux.HandleFunc(msLoginCallbackUrl, servers.AllowMethods(v.msLoginCallback, http.MethodPost))

	mux.ServeHTTP(w, r)
}

// NewAPIServerV1 creates a new APIServerV1. This is a sub-router and should not be used
// as a base router.
func NewAPIServerV1(l *zap.Logger, db *database.DB, azureClient *confidential.Client) *APIServerV1 {
	return &APIServerV1{l, db, azureClient}
}
