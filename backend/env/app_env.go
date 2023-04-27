package env

const (
	appEnv = "APP_ENV"
)

type AppEnv string

const (
	AppEnvLocal       AppEnv = "local"
	AppEnvIntegration        = "integration"
)

// GetAppEnv returns the APP_ENV environment variable.
// Note that this variable is required.
func GetAppEnv() (AppEnv, error) {
	env, err := getRequiredEnv(appEnv)
	return AppEnv(env), err
}
