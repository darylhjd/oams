package v1

import (
	"net/http"

	"go.uber.org/zap"
)

const (
	baseUrl = "/"
	pingUrl = "/ping"
)

type APIServerV1 struct {
	l *zap.Logger
}

func (v *APIServerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc(baseUrl, v.base)
	mux.HandleFunc(pingUrl, v.ping)

	mux.ServeHTTP(w, r)
}

// NewAPIServerV1 creates a new APIServerV1. This is a sub-router and should not be used
// as a base router.
func NewAPIServerV1(l *zap.Logger) *APIServerV1 {
	return &APIServerV1{l: l}
}
