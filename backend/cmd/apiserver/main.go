package main

import (
	"log"

	"github.com/darylhjd/oats/backend/cmd"
	_ "github.com/darylhjd/oats/backend/env/autoloader"
	"github.com/darylhjd/oats/backend/servers/apiserver"
)

func main() {
	apiServer, err := apiserver.New()
	if err != nil {
		log.Fatalf("%s - cannot start %s service: %s\n", cmd.Namespace, apiserver.Namespace, err)
	}

	if err = apiServer.Start(); err != nil {
		apiServer.GetLogger().Error(err.Error())
	}

	if err = apiServer.Stop(); err != nil {
		apiServer.GetLogger().Error(err.Error())
	}
}
