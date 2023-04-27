package apiserver

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/darylhjd/oats/backend/env"
	"github.com/darylhjd/oats/backend/logger"
	"github.com/darylhjd/oats/backend/servers/apiserver/v1"
)

const (
	v1Url = "/api/v1/"
)

// APIServer defines the server structure for the Oats API service.
type APIServer struct {
	l  *zap.Logger
	v1 *v1.APIServerV1
}

// Start the APIServer.
func (s *APIServer) Start() error {
	s.l.Info("apiserver - starting service...")

	port, err := env.GetAPIServerPort()
	if err != nil {
		return err
	}

	s.l.Info("apiserver - service started", zap.String("port", port))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), s)
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	// To add more versions for this URL, simply add another handler for it.
	mux.Handle(v1Url, http.StripPrefix(strings.TrimSuffix(v1Url, "/"), s.v1))

	mux.ServeHTTP(w, r)
}

func (s *APIServer) GetLogger() *zap.Logger {
	return s.l
}

// New creates a new APIServer. Use Start() to start the server.
func New() (*APIServer, error) {
	l, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("apiserver - failed to initialise: %w", err)
	}

	return &APIServer{l, v1.NewAPIServerV1(l)}, nil
}
