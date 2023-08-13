package main

import (
	"context"
	"log"

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

	go func() {
		if err := apiServer.Start(); err != nil {
			log.Fatalf("%s - cannot start apiserver", namespace)
		}
	}()

	go func() {
		if err := webServer.Start(); err != nil {
			log.Fatalf("%s - cannot start webserver", namespace)
		}
	}()

	select {} // Block forever and continue running both servers.
}
