package env

const (
	databaseType     = "DATABASE_TYPE"
	databaseName     = "DATABASE_NAME"
	databaseUser     = "DATABASE_USER"
	databasePassword = "DATABASE_PASSWORD"
	databaseHost     = "DATABASE_HOST"
	databasePort     = "DATABASE_PORT"
)

// GetDatabaseType returns the DATABASE_TYPE environment variable.
// Note that this variable is required.
func GetDatabaseType() (string, error) {
	return getRequiredEnv(databaseType)
}

// GetDatabaseName returns the DATABASE_NAME environment variable.
// Note that this variable is required
func GetDatabaseName() (string, error) {
	return getRequiredEnv(databaseName)
}

// GetDatabaseUser returns the DATABASE_USER environment variable.
// Note that this variable is required
func GetDatabaseUser() (string, error) {
	return getRequiredEnv(databaseUser)
}

// GetDatabasePassword returns the DATABASE_PASSWORD environment variable.
// Note that this variable is required
func GetDatabasePassword() (string, error) {
	return getRequiredEnv(databasePassword)
}

// GetDatabaseHost returns the DATABASE_HOST environment variable.
// Note that this variable is required
func GetDatabaseHost() (string, error) {
	return getRequiredEnv(databaseHost)
}

// GetDatabasePort returns the DATABASE_PORT environment variable.
// Note that this variable is required
func GetDatabasePort() (string, error) {
	return getRequiredEnv(databasePort)
}
