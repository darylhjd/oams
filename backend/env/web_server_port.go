package env

const (
	webServerPort = "WEB_SERVER_PORT"
)

// GetWebServerPort returns the WEB_SERVER_PORT environment variable.
// Note that this variable is required.
func GetWebServerPort() (string, error) {
	return getRequiredEnv(webServerPort)
}
