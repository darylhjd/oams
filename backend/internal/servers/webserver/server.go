package webserver

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/logger"
)

const Namespace = "webserver"

const (
	buildPath = "build/out"

	index    = "index"
	notFound = "404.html"
)

const (
	appUrl = "/"
)

//go:embed build/out/*
var website embed.FS

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

	fileServer, err := fs.Sub(website, buildPath)
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
	// Check if a request refers to a page.
	// It must meet the following condition:
	// 1. No extension in file name.
	path := strings.Trim(r.URL.Path, "/")

	s.l.Debug(fmt.Sprintf("%s - received request", Namespace),
		zap.String("url", r.URL.String()),
		zap.String("cleaned", path))

	if filepath.Ext(path) != "" {
		s.mux.ServeHTTP(w, r)
		return
	}

	// Get appropriate HTML file to serve.
	if path == "" {
		path = index
	}

	file, err := s.app.Open(fmt.Sprintf("%s.html", path))
	if err != nil {
		file, _ = s.app.Open(notFound)
	}

	if _, err = io.Copy(w, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
