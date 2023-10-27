package env

import (
	"fmt"
	"os"
)

const (
	configuration = "CONFIGURATION"
)

type Configuration string

const (
	ConfigurationApiServer = "apiserver"
	ConfigurationDatabase  = "database"
)

// getConfiguration returns the CONFIGURATION environment variable.
func getConfiguration() Configuration {
	return Configuration(os.Getenv(configuration))
}

func verifyApiServerConfiguration() error {
	envs := []string{
		apiServerPort,
		apiServerAzureTenantId,
		apiServerAzureClientId,
		apiServerAzureClientSecret,
		apiServerAzureLoginScope,
		apiServerAzureLoginCallbackUrl,
		webServer,
	}

	if err := checkRequiredEnvs(envs...); err != nil {
		return err
	}

	if err := verifyDatabaseConfiguration(); err != nil {
		return err
	}

	return nil
}

func verifyDatabaseConfiguration() error {
	envs := []string{
		databaseType,
		databaseName,
		databaseUser,
		databasePassword,
		databaseHost,
		databasePort,
		databaseSslMode,
	}

	if err := checkRequiredEnvs(envs...); err != nil {
		return err
	}

	// Check SSL mode and root cert.
	mode := GetDatabaseSSLMode()
	env := GetAppEnv()
	switch mode {
	case databaseSslModeDisable:
		if env == AppEnvStaging {
			return fmt.Errorf("database connection ssl mode disabled for critical environment %q", env)
		}
		fallthrough
	case databaseSslModeVerifyFull:
		return nil
	default:
		return fmt.Errorf("invalid %s value %q", databaseSslMode, mode)
	}
}
