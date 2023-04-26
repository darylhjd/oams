package apiserver

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oats/backend/logger"
	v1 "github.com/darylhjd/oats/backend/servers/apiserver/v1"
)

// APIServer defines the servers structure for the Oats API service.
type APIServer struct {
	L *zap.Logger
}

func (s *APIServer) Start() error {
	s.L.Info("apiserver: starting service...")
	mux := http.NewServeMux()

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", &v1.APIServerV1{L: s.L}))

	s.L.Info("apiserver: service started, serving requests")
	return http.ListenAndServe(":3000", mux)
}

func (s *APIServer) GetLogger() *zap.Logger {
	return s.L
}

// NewAPIServer returns the servers structure for the API service.
func NewAPIServer() (*APIServer, error) {
	l, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("apiserver - failed to initialise: %w", err)
	}

	return &APIServer{l}, nil
}
