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
	case AppEnvLocal, AppEnvStaging, AppEnvProduction:
		return nil
	default:
		return fmt.Errorf("unknown %s value: %s", appEnv, environment)
	}
}

func checkRequiredEnvs(envs ...string) error {
	var errs []error
	for _, env := range envs {
		if checkVarEmpty(os.Getenv(env)) {
			errs = append(errs, fmt.Errorf("%s not set, but is required", env))
		}
	}

	return errors.Join(errs...)
}

func checkVarEmpty(v string) bool {
	if strings.TrimSpace(v) == "" {
		return true
	}

	return false
}
