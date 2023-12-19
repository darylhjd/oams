package env

import "os"

const (
	appEnv = "APP_ENV"
)

type AppEnv string

const (
	AppEnvLocal      AppEnv = "local"
	AppEnvStaging           = "staging"
	AppEnvProduction        = "production"
)

// GetAppEnv returns the APP_ENV environment variable.
func GetAppEnv() AppEnv {
	return AppEnv(os.Getenv(appEnv))
}
