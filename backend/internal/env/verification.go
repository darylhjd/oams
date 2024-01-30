package env

import (
	"errors"
	"fmt"
	"os"
)

func verifyEnvironment() error {
	if err := checkEnvsNotEmpty(appEnv); err != nil {
		return err
	}

	switch environment := GetAppEnv(); environment {
	case AppEnvLocal, AppEnvStaging, AppEnvProduction:
		return nil
	default:
		return fmt.Errorf("unknown %s value: %s", appEnv, environment)
	}
}

func verifyConfiguration() error {
	var envs []string
	switch GetConfiguration() {
	case ConfAPIServer:
		envs = []string{
			apiServerPort,
			apiServerAzureTenantId,
			apiServerAzureClientId,
			apiServerAzureClientSecret,
			apiServerAzureLoginScope,
			webServer,
			databaseType,
			databaseName,
			databaseUser,
			databasePassword,
			databaseHost,
			databasePort,
			databaseSslMode,
		}
	case ConfIntervention:
		envs = []string{
			databaseType,
			databaseName,
			databaseUser,
			databasePassword,
			databaseHost,
			databasePort,
			databaseSslMode,
			azureEmailEndpoint,
			azureEmailAccessKey,
			azureEmailSenderAddress,
		}
	default:
		envs = []string{
			apiServerPort,
			apiServerAzureTenantId,
			apiServerAzureClientId,
			apiServerAzureClientSecret,
			apiServerAzureLoginScope,
			webServer,
			databaseType,
			databaseName,
			databaseUser,
			databasePassword,
			databaseHost,
			databasePort,
			databaseSslMode,
			azureEmailEndpoint,
			azureEmailAccessKey,
			azureEmailSenderAddress,
		}
	}

	if err := checkEnvsNotEmpty(envs...); err != nil {
		return err
	}

	if err := verifyDatabaseSsl(); err != nil {
		return err
	}

	return nil
}

func verifyDatabaseSsl() error {
	// Check SSL mode and root cert.
	mode := GetDatabaseSSLMode()
	env := GetAppEnv()
	switch mode {
	case databaseSslModeVerifyFull:
		return checkEnvsNotEmpty(databaseSslRootCertLoc)
	case databaseSslModeDisable:
		if env == AppEnvLocal {
			return nil
		}

		return fmt.Errorf("database connection ssl mode disabled for critical environment %q", env)
	default:
		return fmt.Errorf("invalid %s value %q", databaseSslMode, mode)
	}
}

func checkEnvsNotEmpty(envs ...string) error {
	var errs []error
	for _, env := range envs {
		if _, ok := os.LookupEnv(env); !ok {
			errs = append(errs, fmt.Errorf("%s not set, but is required", env))
		}
	}

	return errors.Join(errs...)
}
