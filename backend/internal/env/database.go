package env

import (
	"os"
)

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
func GetDatabaseType() string {
	return os.Getenv(databaseType)
}

// GetDatabaseName returns the DATABASE_NAME environment variable.
func GetDatabaseName() string {
	return os.Getenv(databaseName)
}

// GetDatabaseUser returns the DATABASE_USER environment variable.
func GetDatabaseUser() string {
	return os.Getenv(databaseUser)
}

// GetDatabasePassword returns the DATABASE_PASSWORD environment variable.
func GetDatabasePassword() string {
	return os.Getenv(databasePassword)
}

// GetDatabaseHost returns the DATABASE_HOST environment variable.
func GetDatabaseHost() string {
	return os.Getenv(databaseHost)
}

// GetDatabasePort returns the DATABASE_PORT environment variable.
func GetDatabasePort() string {
	return os.Getenv(databasePort)
}

// GetDatabaseSSLMode returns the DATABASE_SSL_MODE environment variable.
func GetDatabaseSSLMode() string {
	return os.Getenv(databaseSslMode)
}

// GetDatabaseSSLRootCertLocation returns the DATABASE_SSL_ROOT_CERT_LOC environment variable.
func GetDatabaseSSLRootCertLocation() string {
	return os.Getenv(databaseSslRootCertLoc)
}
