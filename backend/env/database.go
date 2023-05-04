package env

import "fmt"

const (
	databaseType           = "DATABASE_TYPE"
	databaseName           = "DATABASE_NAME"
	databaseUser           = "DATABASE_USER"
	databasePassword       = "DATABASE_PASSWORD"
	databaseHost           = "DATABASE_HOST"
	databasePort           = "DATABASE_PORT"
	databaseSslMode        = "DATABASE_SSL_MODE"
	databaseSslRootCertLoc = "DATABASE_SSL_ROOT_CERT_LOC"
)

const (
	databaseSslModeVerifyFull = "verify-full"
	databaseSslModeDisable    = "disable"
)

// GetDatabaseType returns the DATABASE_TYPE environment variable.
// Note that this variable is required.
func GetDatabaseType() (string, error) {
	return getRequiredEnv(databaseType)
}

// GetDatabaseName returns the DATABASE_NAME environment variable.
// Note that this variable is required.
func GetDatabaseName() (string, error) {
	return getRequiredEnv(databaseName)
}

// GetDatabaseUser returns the DATABASE_USER environment variable.
// Note that this variable is required.
func GetDatabaseUser() (string, error) {
	return getRequiredEnv(databaseUser)
}

// GetDatabasePassword returns the DATABASE_PASSWORD environment variable.
// Note that this variable is required.
func GetDatabasePassword() (string, error) {
	return getRequiredEnv(databasePassword)
}

// GetDatabaseHost returns the DATABASE_HOST environment variable.
// Note that this variable is required.
func GetDatabaseHost() (string, error) {
	return getRequiredEnv(databaseHost)
}

// GetDatabasePort returns the DATABASE_PORT environment variable.
// Note that this variable is required.
func GetDatabasePort() (string, error) {
	return getRequiredEnv(databasePort)
}

// GetDatabaseSSLMode returns the DATABASE_SSL_MODE environment variable.
// Note that this variable is required.
func GetDatabaseSSLMode() (string, error) {
	mode, err := getRequiredEnv(databaseSslMode)
	if err != nil {
		return "", err
	}

	env, err := GetAppEnv()
	if err != nil {
		return "", err
	}

	switch mode {
	case databaseSslModeDisable:
		if env == AppEnvStaging {
			return "", fmt.Errorf("%s - database connection ssl mode disabled for critical environment %q", namespace, env)
		}
		fallthrough
	case databaseSslModeVerifyFull:
		return mode, nil
	default:
		return "", fmt.Errorf("%s - invalid %s value %q", namespace, databaseSslMode, mode)
	}
}

// GetDatabaseSSLRootCertLocation returns the DATABASE_SSL_ROOT_CERT_LOC environment variable.
// Note that this variable is required depending on whether SSL mode is enabled or disabled.
func GetDatabaseSSLRootCertLocation() (string, error) {
	sslMode, err := GetDatabaseSSLMode()
	if err != nil {
		return "", err
	}

	if sslMode == databaseSslModeDisable {
		return "", nil
	}

	return getRequiredEnv(databaseSslRootCertLoc)
}
