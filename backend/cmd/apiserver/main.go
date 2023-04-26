package main

import (
	"log"

	"github.com/darylhjd/oats/backend/servers/apiserver"
)

func main() {
	apiServer, err := apiserver.NewAPIServer()
	if err != nil {
		log.Fatalf("cmd - cannot start apiserver service: %s", err)
	}

	if err := apiServer.Start(); err != nil {
		apiServer.L.Error(err.Error())
	}
}
