package logger

import (
	"log"

	"go.uber.org/zap"

	"github.com/darylhjd/oats/backend/env"
)

// NewLogger returns a logger with the specified logging LogLevel.
func NewLogger() (*zap.Logger, error) {
	logLevel := env.GetLogLevel()

	switch logLevel {
	case env.LogLevelProd:
		return zap.NewProduction()
	case env.LogLevelDev:
		return zap.NewDevelopment()
	default:
		log.Printf("%q not a valid log level, using default development logging.\n", logLevel)
		return zap.NewDevelopment()
	}
}
