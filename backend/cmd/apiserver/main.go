package main

import (
	"context"
	"log"

	"github.com/darylhjd/oams/backend/internal/servers/apiserver"
)

func main() {
	apiServer, err := apiserver.New(context.Background())
	if err != nil {
		log.Fatalf("%s - cannot start service: %s\n", apiserver.Namespace, err)
	}
	defer func() {
		if err = apiServer.Stop(); err != nil {
			log.Fatalf("%s - %s", apiserver.Namespace, err)
		}
	}()

	if err = apiServer.Start(); err != nil {
		log.Fatalf("%s - %s", apiserver.Namespace, err)
	}
}
