package env

import "os"

const (
	webServer = "WEB_SERVER"
)

// GetWebServer returns the WEB_SERVER environment variable.
func GetWebServer() string {
	return os.Getenv(webServer)
}
