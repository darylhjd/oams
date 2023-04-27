package autoloader

import (
	"log"

	"github.com/darylhjd/oats/backend/env"
)

var loaded bool

func init() {
	if !loaded {
		log.Println("Loaded envs")
		env.MustLoad()
		loaded = !loaded
	}
}
