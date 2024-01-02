package main

import (
	"context"
	"log"

	"github.com/darylhjd/oams/backend/internal/intervention"
)

func main() {
	ctx := context.Background()

	service, err := intervention.New(ctx)
	if err != nil {
		log.Fatalf("%s - cannot start service: %s", intervention.Namespace, err)
	}
	defer func() {
		if err = service.Stop(); err != nil {
			log.Printf("%s - could not gracefully stop service: %s", intervention.Namespace, err)
		}
	}()

	if err = service.Run(); err != nil {
		log.Printf("%s - error running service: %s", intervention.Namespace, err)
	}
}
