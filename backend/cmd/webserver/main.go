package main

import (
	"log"

	"github.com/darylhjd/oams/backend/internal/servers/webserver"
)

func main() {
	webServer, err := webserver.New()
	if err != nil {
		log.Fatalf("%s - cannot start service: %s\n", webserver.Namespace, err)
	}

	if err = webServer.Start(); err != nil {
		log.Fatalf("%s - %s", webserver.Namespace, err)
	}
}
