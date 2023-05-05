package env

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
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

func verifyEnvironment() error {
	environment := GetAppEnv()
	if err := checkVarEmpty(string(environment)); err != nil {
		return err
	}

	var err error
	baseErr := fmt.Errorf("%s - environment verification failed", namespace)

	switch environment {
	case AppEnvLocal:
		err = godotenv.Load(GetEnvLoc())
	case AppEnvIntegration, AppEnvStaging:
		return nil
	default:
		err = fmt.Errorf("%s - unknown %s value: %s", namespace, appEnv, environment)
	}

	if err != nil {
		return errors.Join(baseErr, err)
	}

	return nil
}

func verifyConfiguration() error {
	switch GetConfiguration() {
	case ConfigurationApiServer:
		return verifyApiServerConfiguration()
	case ConfigurationWebServer:
		return verifyWebServerConfiguration()
	case ConfigurationDatabase:
		return verifyDatabaseConfiguration()
	default:
		// If no configuration is specified, all configurations are checked.
		return errors.Join(
			verifyApiServerConfiguration(),
			verifyWebServerConfiguration(),
			verifyDatabaseConfiguration())
	}
}

func checkRequiredEnvs(envs ...string) error {
	var errs []error
	for _, env := range envs {
		if err := checkVarEmpty(os.Getenv(env)); err != nil {
			errs = append(errs, fmt.Errorf("%s not set, but is required", env))
		}
	}

	return errors.Join(errs...)
}

func checkVarEmpty(v string) error {
	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("%s - empty variable", namespace)
	}

	return nil
}
