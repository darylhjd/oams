package apiserver

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/darylhjd/oats/backend/database"
	"github.com/darylhjd/oats/backend/env"
	"github.com/darylhjd/oats/backend/logger"
	"github.com/darylhjd/oats/backend/servers/apiserver/v1"
)

const (
	Namespace = "apiserver"

	v1Url = "/api/v1/"
)

// APIServer defines the server structure for the Oats API service.
type APIServer struct {
	l  *zap.Logger
	db *database.DB

	v1 *v1.APIServerV1
}

// Start the APIServer.
func (s *APIServer) Start() error {
	s.l.Info(fmt.Sprintf("%s - starting service...", Namespace))

	port, err := env.GetAPIServerPort()
	if err != nil {
		return err
	}

	s.l.Info(fmt.Sprintf("%s - service started", Namespace), zap.String("port", port))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), s)
}

// Stop closes any external connections (e.g. database) and stops the server gracefully.
func (s *APIServer) Stop() error {
	return s.db.Close()
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
		return nil, fmt.Errorf("%s - failed to initialise: %w", Namespace, err)
	}

	db, err := database.Connect()
	if err != nil {
		return nil, fmt.Errorf("%s - could not connect to database: %w", Namespace, err)
	}

	return &APIServer{l, db, v1.NewAPIServerV1(l, db)}, nil
}
