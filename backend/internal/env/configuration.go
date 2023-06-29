package env

import (
	"errors"
	"fmt"
	"os"
)

const (
	configuration = "CONFIGURATION"
)

type Configuration string

const (
	ConfigurationApiServer = "apiserver"
	ConfigurationWebServer = "webserver"
	ConfigurationDatabase  = "database"
)

// GetConfiguration returns the CONFIGURATION environment variable.
func GetConfiguration() Configuration {
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
		webServerHost,
		webServerPort,
	}

	baseErr := fmt.Errorf("%s - configuration verification failed for %s", namespace, ConfigurationApiServer)
	if err := checkRequiredEnvs(envs...); err != nil {
		return errors.Join(baseErr, err)
	}

	if err := verifyDatabaseConfiguration(); err != nil {
		return errors.Join(baseErr, err)
	}

	return nil
}

func verifyWebServerConfiguration() error {
	envs := []string{
		webServerPort,
	}

	baseErr := fmt.Errorf("%s - configuration verification failed for %s", namespace, ConfigurationWebServer)
	if err := checkRequiredEnvs(envs...); err != nil {
		return errors.Join(baseErr, err)
	}

	if err := verifyDatabaseConfiguration(); err != nil {
		return errors.Join(baseErr, err)
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

	baseErr := fmt.Errorf("%s - configuration verification failed for %s", namespace, ConfigurationDatabase)
	if err := checkRequiredEnvs(envs...); err != nil {
		return errors.Join(baseErr, err)
	}

	// Check SSL mode and root cert.
	mode := GetDatabaseSSLMode()
	env := GetAppEnv()
	switch mode {
	case databaseSslModeDisable:
		if env == AppEnvStaging {
			disabledErr := fmt.Errorf("database connection ssl mode disabled for critical environment %q", env)
			return errors.Join(baseErr, disabledErr)
		}
		fallthrough
	case databaseSslModeVerifyFull:
		return nil
	default:
		invalidErr := fmt.Errorf("invalid %s value %q", databaseSslMode, mode)
		return errors.Join(baseErr, invalidErr)
	}
}
