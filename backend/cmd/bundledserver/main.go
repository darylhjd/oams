package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver"
	"github.com/darylhjd/oams/backend/internal/servers/webserver"
)

const (
	namespace = "bundledserver"
)

func main() {
	apiServer, err := apiserver.New(context.Background())
	if err != nil {
		log.Fatalf("%s - cannot create apiserver", namespace)
	}

	webServer, err := webserver.New()
	if err != nil {
		log.Fatalf("%s - cannot create webserver", namespace)
	}

	b := BundledServer{apiServer, webServer}

	// Default to the APIServer port.
	if err = http.ListenAndServe(fmt.Sprintf(":%s", env.GetAPIServerPort()), b); err != nil {
		b.apiServer.Stop()
		log.Fatalf("%s - %s", namespace, err)
	}
}

// BundledServer allows us to run both the APIServer and the WebServer in one service.
type BundledServer struct {
	apiServer *apiserver.APIServer
	webServer *webserver.WebServer
}

// ServeHTTP allows us to fulfill the http.Handler interface.
func (b BundledServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, apiserver.Url) {
		b.apiServer.ServeHTTP(w, r)
	} else {
		b.webServer.ServeHTTP(w, r)
	}
}
