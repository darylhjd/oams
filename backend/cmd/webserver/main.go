package main

import (
	"log"

	"github.com/darylhjd/oats/backend/env"
	"github.com/darylhjd/oats/backend/servers/webserver"
)

func main() {
	env.MustLoad()

	webServer, err := webserver.New()
	if err != nil {
		log.Fatalf("cmd - cannot start webserver service: %s", err)
	}

	if err = webServer.Start(); err != nil {
		webServer.GetLogger().Error(err.Error())
	}
}
