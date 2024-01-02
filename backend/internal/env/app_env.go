package env

import "os"

const (
	appEnv = "APP_ENV"
)

type AppEnv string

const (
	AppEnvLocal      AppEnv = "local"
	AppEnvStaging    AppEnv = "staging"
	AppEnvProduction AppEnv = "production"
)

// GetAppEnv returns the APP_ENV environment variable.
func GetAppEnv() AppEnv {
	return AppEnv(os.Getenv(appEnv))
}
