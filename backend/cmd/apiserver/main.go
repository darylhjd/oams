package main

import (
	"context"
	"log"

	"github.com/darylhjd/oams/backend/servers/apiserver"
)

func main() {
	apiServer, err := apiserver.New(context.Background())
	if err != nil {
		log.Fatalf("%s - cannot start service: %s\n", apiserver.Namespace, err)
	}

	if err = apiServer.Start(); err != nil {
		apiServer.GetLogger().Error(err.Error())
	}

	apiServer.Stop()
}
