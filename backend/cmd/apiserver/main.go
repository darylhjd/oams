package main

import (
	"log"

	_ "github.com/darylhjd/oats/backend/env/autoloader"
	"github.com/darylhjd/oats/backend/servers/apiserver"
)

func main() {
	apiServer, err := apiserver.New()
	if err != nil {
		log.Fatalf("cmd - cannot start apiserver service: %s", err)
	}

	if err = apiServer.Start(); err != nil {
		apiServer.GetLogger().Error(err.Error())
	}
}
