package env

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const namespace = "env"

// MustLoad loads environment variables depending on the environment.
func MustLoad() {
	environment, err := GetAppEnv()
	if err != nil {
		log.Fatal(err)
	}

	switch environment {
	case AppEnvLocal:
		err = godotenv.Load(GetEnvLoc())
	case AppEnvIntegration, AppEnvStaging:
		return
	default:
		err = fmt.Errorf("%s - unknown %s value: %s", namespace, appEnv, environment)
	}

	if err != nil {
		log.Fatal(fmt.Errorf("%s - could not load: %w", namespace, err))
	}
}

func getRequiredEnv(env string) (string, error) {
	e := os.Getenv(env)
	if strings.TrimSpace(e) == "" {
		return "", fmt.Errorf("%s - %s not set, but is compulsory", namespace, env)
	}

	return e, nil
}
