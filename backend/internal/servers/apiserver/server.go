package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/cors"
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/logger"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/v1"
)

const (
	Namespace = "apiserver"
	Url       = "/api/"
)

// APIServer defines the server structure for the OAMS API service.
type APIServer struct {
	l   *zap.Logger
	db  *database.DB
	mux *http.ServeMux

	v1 *v1.APIServerV1
}

// New creates a new APIServer. Use Start() to start the server.
func New(ctx context.Context) (*APIServer, error) {
	l, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("%s - failed to initialise: %w", Namespace, err)
	}

	db, err := database.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s - could not connect to database: %w", Namespace, err)
	}

	azureAuthenticator, err := oauth2.NewAzureAuthenticator()
	if err != nil {
		return nil, fmt.Errorf("%s - could not create azure authenticator: %w", Namespace, err)
	}

	server := &APIServer{l, db, http.NewServeMux(), v1.New(l, db, azureAuthenticator)}
	server.registerHandlers()

	return server, nil
}

// Start the APIServer.
func (s *APIServer) Start() error {
	s.l.Info(fmt.Sprintf("%s - starting service...", Namespace))

	// Set up CORS.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{env.GetWebServer()},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization, Content-Type"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
	})

	port := env.GetAPIServerPort()
	s.l.Info(fmt.Sprintf("%s - service started", Namespace), zap.String("port", port))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), c.Handler(s))
}

func (s *APIServer) registerHandlers() {
	s.mux.HandleFunc(Url, s.base)

	// To add more versions for this URL, simply add another handler for it.
	s.mux.Handle(v1.Url, http.StripPrefix(strings.TrimSuffix(v1.Url, "/"), s.v1))
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.l.Debug(fmt.Sprintf("%s - deciding version", Namespace),
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)

	s.mux.ServeHTTP(w, r)
}

func (s *APIServer) base(w http.ResponseWriter, _ *http.Request) {
	resp := struct {
		Message string `json:"message"`
	}{
		"malformed url path",
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, string(respBytes), http.StatusNotFound)
}

// Stop closes any external connections (e.g. database) and stops the server gracefully.
func (s *APIServer) Stop() error {
	return s.db.Close()
}
