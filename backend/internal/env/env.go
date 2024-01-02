package env

import (
	"fmt"
	"log"
)

const namespace = "env"

// MustLoad loads environment variables depending on the environment.
func MustLoad() {
	if err := verifyEnvironment(); err != nil {
		log.Fatal(fmt.Errorf("%s - environment check failed: %w", namespace, err))
	}

	if err := verifyConfiguration(); err != nil {
		log.Fatal(fmt.Errorf("%s - configuration check failed: %w", namespace, err))
	}
}
