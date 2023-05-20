package webserver

import (
	"fmt"
	"io/fs"
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/frontend"

	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/logger"
)

const Namespace = "webserver"

const (
	buildPath = "build/web"
)

const (
	appUrl = "/"
)

// WebServer defines the server structure for the OAMS Web Server.
type WebServer struct {
	l   *zap.Logger
	app fs.FS
	mux *http.ServeMux
}

// New creates a new WebServer. Use Start() to start the server.
func New() (*WebServer, error) {
	l, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("%s - failed to initialise logger: %w", Namespace, err)
	}

	fileServer, err := fs.Sub(frontend.Website, buildPath)
	if err != nil {
		return nil, fmt.Errorf("%s - failed to initialise: %w", Namespace, err)
	}

	server := &WebServer{l, fileServer, http.NewServeMux()}
	server.registerHandlers()

	return server, nil
}

// Start the WebServer.
func (s *WebServer) Start() error {
	s.l.Info(fmt.Sprintf("%s - starting service...", Namespace))

	port := env.GetWebServerPort()
	s.l.Info(fmt.Sprintf("%s - service started on port %s", Namespace, port))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), s)
}

func (s *WebServer) registerHandlers() {
	// Add your paths here.
	s.mux.Handle(appUrl, http.FileServer(http.FS(s.app)))
}

func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.l.Debug(fmt.Sprintf("%s - received request", Namespace),
		zap.String("url", r.URL.String()))

	s.mux.ServeHTTP(w, r)
}

func (s *WebServer) GetLogger() *zap.Logger {
	return s.l
}
