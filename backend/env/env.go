package env

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// MustLoad loads environment variables depending on the environment.
func MustLoad() {
	environment, err := GetAppEnv()
	if err != nil {
		log.Fatal(err)
	}

	switch environment {
	case AppEnvLocal:
		err = godotenv.Load(GetEnvLoc())
	default:
		err = fmt.Errorf("env - unknown %s value: %s", appEnv, environment)
	}

	if err != nil {
		log.Fatal(fmt.Errorf("env - could not load: %w", err))
	}
}

func getRequiredEnv(env string) (string, error) {
	e := os.Getenv(env)
	if strings.TrimSpace(e) == "" {
		return "", fmt.Errorf("env - %s not set, but is compulsory", env)
	}

	return e, nil
}
