package env

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// MustLoad loads environment variables.
func MustLoad() {
	env, err := GetAppEnv()
	if err != nil {
		log.Fatal(err)
	}

	switch env {
	case AppEnvLocal:
		err = godotenv.Load(".env.local")
	default:
		err = fmt.Errorf("env - unknown %s value: %s", appEnv, env)
	}

	if err != nil {
		log.Fatal(err)
	}
}
