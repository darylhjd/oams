package main

import (
	"log"

	_ "github.com/darylhjd/oams/backend/env/autoloader"
	"github.com/darylhjd/oams/backend/servers/apiserver"
)

func main() {
	apiServer, err := apiserver.New()
	if err != nil {
		log.Fatalf("%s - cannot start service: %s\n", apiserver.Namespace, err)
	}

	if err = apiServer.Start(); err != nil {
		apiServer.GetLogger().Error(err.Error())
	}

	if err = apiServer.Stop(); err != nil {
		apiServer.GetLogger().Error(err.Error())
	}
}
