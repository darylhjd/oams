package env

import (
	"fmt"
	"os"
	"strings"
)

const (
	apiServerPort = "API_SERVER_PORT"
)

// GetAPIServerPort returns the API_SERVER_PORT environment variable.
// Note that this variable is required.
func GetAPIServerPort() (string, error) {
	port := os.Getenv(apiServerPort)
	if strings.TrimSpace(port) == "" {
		return "", fmt.Errorf("env - %s not set, but is compulsory", apiServerPort)
	}

	return port, nil
}
