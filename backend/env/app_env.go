package env

import "os"

const (
	appEnv = "APP_ENV"
)

type AppEnv string

const (
	AppEnvLocal       AppEnv = "local"
	AppEnvIntegration        = "integration"
	AppEnvStaging            = "staging"
)

// GetAppEnv returns the APP_ENV environment variable.
func GetAppEnv() AppEnv {
	return AppEnv(os.Getenv(appEnv))
}
