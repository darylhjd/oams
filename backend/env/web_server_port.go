package env

import "os"

const (
	webServerPort = "WEB_SERVER_PORT"
)

// GetWebServerPort returns the WEB_SERVER_PORT environment variable.
func GetWebServerPort() string {
	return os.Getenv(webServerPort)
}
