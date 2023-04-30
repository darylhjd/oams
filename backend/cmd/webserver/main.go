package main

import (
	"log"

	_ "github.com/darylhjd/oats/backend/env/autoloader"
	"github.com/darylhjd/oats/backend/servers/webserver"
)

func main() {
	webServer, err := webserver.New()
	if err != nil {
		log.Fatalf("%s - cannot start service: %s\n", webserver.Namespace, err)
	}

	if err = webServer.Start(); err != nil {
		webServer.GetLogger().Error(err.Error())
	}
}
