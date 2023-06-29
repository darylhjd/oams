package env

import "os"

const (
	webServerHost = "WEB_SERVER_HOST"
	webServerPort = "WEB_SERVER_PORT"
)

// GetWebServerHost returns the WEB_SERVER_HOST environment variable.
func GetWebServerHost() string {
	return os.Getenv(webServerHost)
}

// GetWebServerPort returns the WEB_SERVER_PORT environment variable.
func GetWebServerPort() string {
	return os.Getenv(webServerPort)
}
