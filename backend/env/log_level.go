package env

import "os"

const (
	logLevel = "LOG_LEVEL"
)

type LogLevel string

const (
	LogLevelDev  LogLevel = "0"
	LogLevelProd          = "1"
)

// GetLogLevel returns the LogLevel specified in the environment variable.
func GetLogLevel() LogLevel {
	return LogLevel(os.Getenv(logLevel))
}
