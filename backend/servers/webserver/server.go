package webserver

import (
	"fmt"
	"io/fs"
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oats/frontend"

	"github.com/darylhjd/oats/backend/env"
	"github.com/darylhjd/oats/backend/logger"
)

const (
	homeUrl = "/"
)

// WebServer defines the server structure for the Oats Web Server.
type WebServer struct {
	l   *zap.Logger
	app fs.FS
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
	s.l.Debug("webserver - received request",
		zap.String("url", r.URL.String()))

	mux := http.NewServeMux()

	// Add your paths here.
	mux.Handle(homeUrl, http.FileServer(http.FS(s.app)))

	mux.ServeHTTP(w, r)
}

func (s *WebServer) GetLogger() *zap.Logger {
	return s.l
}

// New creates a new WebServer. Use Start() to start the server.
func New() (*WebServer, error) {
	l, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("webserver - failed to initialise logger: %w", err)
	}

	server, err := fs.Sub(frontend.Website, "build")
	if err != nil {
		return nil, fmt.Errorf("webserver - failed to initialise: %w", err)
	}

	return &WebServer{l, server}, nil
}
