package autoloader

import (
	"github.com/darylhjd/oats/backend/env"
)

var loaded bool

func init() {
	if !loaded {
		env.MustLoad()
		loaded = !loaded
	}
}
