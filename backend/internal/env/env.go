package env

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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
	if err := checkRequiredEnvs(appEnv); err != nil {
		return err
	}

	switch environment := GetAppEnv(); environment {
	case AppEnvLocal, AppEnvIntegration, AppEnvStaging:
		return nil
	default:
		return fmt.Errorf("unknown %s value: %s", appEnv, environment)
	}
}

func verifyConfiguration() error {
	switch getConfiguration() {
	case ConfigurationApiServer:
		return verifyApiServerConfiguration()
	case ConfigurationDatabase:
		return verifyDatabaseConfiguration()
	default:
		// If no configuration is specified, all configurations are checked.
		return errors.Join(
			verifyApiServerConfiguration(),
			verifyDatabaseConfiguration(),
		)
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
		return errors.New("empty variable")
	}

	return nil
}
