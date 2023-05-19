package webserver

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/frontend"

	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/logger"
)

const Namespace = "webserver"

const (
	buildPath = "build"
	indexPath = "index.html"
)

const (
	appUrl = "/"
)

// WebServer defines the server structure for the OAMS Web Server.
type WebServer struct {
	l   *zap.Logger
	app fs.FS
}

// New creates a new WebServer. Use Start() to start the server.
func New() (*WebServer, error) {
	l, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("%s - failed to initialise logger: %w", Namespace, err)
	}

	server, err := fs.Sub(frontend.Website, "build")
	if err != nil {
		return nil, fmt.Errorf("%s - failed to initialise: %w", Namespace, err)
	}

	return &WebServer{l, server}, nil
}

// Start the WebServer.
func (s *WebServer) Start() error {
	s.l.Info(fmt.Sprintf("%s - starting service...", Namespace))

	port := env.GetWebServerPort()
	s.l.Info(fmt.Sprintf("%s - service started on port %s", Namespace, port))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), s)
}

func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.l.Debug(fmt.Sprintf("%s - received request", Namespace), zap.String("url", r.URL.String()))

	path := r.URL.Path
	if path == appUrl {
		path = indexPath
	}

	path = strings.TrimPrefix(path, appUrl)

	s.l.Debug(fmt.Sprintf("%s - received request for filepath", Namespace), zap.String("filepath", path))

	// We try to get the file provided from the path.
	file, err := s.app.Open(path)
	if err == nil {
		// No error, file exists so we serve it.
		s.l.Debug(fmt.Sprintf("%s - file exists, serving it", Namespace))
		if _, err = io.Copy(w, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if !errors.Is(err, fs.ErrNotExist) {
		s.l.Debug(fmt.Sprintf("%s - error getting file", Namespace))
		// If there was an error reading the file that is not that it does not exist.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.l.Debug(fmt.Sprintf("%s - filepath does not exist", Namespace))

	// If the file does not exist, we serve the index HTML if that is the file type we expect.
	// Else, we should return 404.
	if ext := filepath.Ext(path); ext != "" && ext != "html" {
		http.NotFound(w, r)
		return
	}

	index, err := s.app.Open(indexPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = io.Copy(w, index); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *WebServer) GetLogger() *zap.Logger {
	return s.l
}
