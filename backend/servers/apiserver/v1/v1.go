package v1

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/database"
	"github.com/darylhjd/oams/backend/servers"
)

const (
	namespace = "apiserver/v1"

	baseUrl = "/"
	pingUrl = "/ping"
)

type APIServerV1 struct {
	l  *zap.Logger
	db *database.DB
}

func (v *APIServerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc(baseUrl, servers.AllowMethods(v.base, http.MethodGet))
	mux.HandleFunc(pingUrl, servers.AllowMethods(v.ping, http.MethodGet))

	mux.ServeHTTP(w, r)
}

// NewAPIServerV1 creates a new APIServerV1. This is a sub-router and should not be used
// as a base router.
func NewAPIServerV1(l *zap.Logger, db *database.DB) *APIServerV1 {
	return &APIServerV1{l, db}
}
