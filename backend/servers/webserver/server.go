package webserver

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oats/backend/env"
	"github.com/darylhjd/oats/backend/logger"
)

// WebServer defines the server structure for the Oats Web Server.
type WebServer struct {
	l *zap.Logger
}

// Start the WebServer.
func (s *WebServer) Start() error {
	s.l.Info("webserver - starting service...")

	port, err := env.GetWebServerPort()
	if err != nil {
		return err
	}

	s.l.Info(fmt.Sprintf("webserver - service started on port %s", port))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), s)
}

func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	// Add your paths here.
	mux.HandleFunc(homeUrl, s.home)

	mux.ServeHTTP(w, r)
}

func (s *WebServer) GetLogger() *zap.Logger {
	return s.l
}

// NewWebServer creates a new WebServer. Use Start() to start the server.
func NewWebServer() (*WebServer, error) {
	l, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("webserver - failed to initialise: %w", err)
	}

	return &WebServer{l}, nil
}
