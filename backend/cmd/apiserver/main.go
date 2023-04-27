package main

import (
	"log"

	"github.com/darylhjd/oats/backend/env"
	_ "github.com/darylhjd/oats/backend/env/autoloader"
	"github.com/darylhjd/oats/backend/servers/apiserver"
)

func main() {
	env.MustLoad()

	apiServer, err := apiserver.New()
	if err != nil {
		log.Fatalf("cmd - cannot start apiserver service: %s", err)
	}

	if err = apiServer.Start(); err != nil {
		apiServer.GetLogger().Error(err.Error())
	}
}
