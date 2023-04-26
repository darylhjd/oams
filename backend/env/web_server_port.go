package env

import (
	"fmt"
	"os"
	"strings"
)

const (
	webServerPort = "WEB_SERVER_PORT"
)

// GetWebServerPort returns the WEB_SERVER_PORT environment variable.
// Note that this variable is required.
func GetWebServerPort() (string, error) {
	port := os.Getenv(webServerPort)
	if strings.TrimSpace(port) == "" {
		return "", fmt.Errorf("env - %s not set, but is compulsory", webServerPort)
	}

	return port, nil
}
