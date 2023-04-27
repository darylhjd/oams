package env

const (
	apiServerPort = "API_SERVER_PORT"
)

// GetAPIServerPort returns the API_SERVER_PORT environment variable.
// Note that this variable is required.
func GetAPIServerPort() (string, error) {
	return getRequiredEnv(apiServerPort)
}
