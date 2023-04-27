package env

import "os"

const (
	envLoc = "ENV_LOC"
)

// GetEnvLoc returns the ENV_LOC environment variable.
func GetEnvLoc() string {
	return os.Getenv(envLoc)
}
