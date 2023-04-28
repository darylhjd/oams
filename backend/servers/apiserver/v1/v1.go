package v1

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oats/backend/database"
)

const (
	baseUrl = "/"
	pingUrl = "/ping"
)

type APIServerV1 struct {
	l  *zap.Logger
	db *database.Queries
}

func (v *APIServerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc(baseUrl, v.base)
	mux.HandleFunc(pingUrl, v.ping)

	mux.ServeHTTP(w, r)
}

// NewAPIServerV1 creates a new APIServerV1. This is a sub-router and should not be used
// as a base router.
func NewAPIServerV1(l *zap.Logger, db *database.Queries) *APIServerV1 {
	return &APIServerV1{l, db}
}
