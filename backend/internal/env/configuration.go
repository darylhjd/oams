package env

import (
	"fmt"
)

func verifyConfiguration() error {
	envs := []string{
		apiServerPort,
		apiServerAzureTenantId,
		apiServerAzureClientId,
		webServer,
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
	case databaseSslModeVerifyFull:
		return checkRequiredEnvs(databaseSslRootCertLoc)
	case databaseSslModeDisable:
		if env == AppEnvLocal {
			return nil
		}

		return fmt.Errorf("database connection ssl mode disabled for critical environment %q", env)
	default:
		return fmt.Errorf("invalid %s value %q", databaseSslMode, mode)
	}
}
