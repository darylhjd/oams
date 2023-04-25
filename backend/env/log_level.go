package env

import "os"

const (
	LogLevel Var = "LOG_LEVEL"
)

// Level defines the levels of logging that is supported by Oats.
type Level string

const (
	LogLevelDev  Level = "0"
	LogLevelProd       = "1"
)

// GetLogLevel returns the Level of logging specified in the environment variable LogLevel.
func GetLogLevel() Level {
	return Level(os.Getenv(string(LogLevel)))
}
