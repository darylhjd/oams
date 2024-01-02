package env

import "os"

const (
	logLevel = "LOG_LEVEL"
)

type LogLevel string

const (
	LogLevelProd  LogLevel = "0"
	LogLevelDebug LogLevel = "1"
)

// GetLogLevel returns the LogLevel specified in the environment variable.
func GetLogLevel() LogLevel {
	return LogLevel(os.Getenv(logLevel))
}
