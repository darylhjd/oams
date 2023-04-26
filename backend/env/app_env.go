package env

import (
	"fmt"
	"os"
)

const (
	appEnv = "APP_ENV"
)

type AppEnv string

const (
	AppEnvLocal AppEnv = "local"
)

// GetAppEnv returns the APP_ENV environment variable.
// Note that this variable is required.
func GetAppEnv() (AppEnv, error) {
	env, ok := os.LookupEnv(appEnv)
	if !ok {
		return "", fmt.Errorf("env - %s not set, but is compulsory", appEnv)
	}

	return AppEnv(env), nil
}
