package logger

import (
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/env"
)

// NewLogger returns a logger with the specified logging LogLevel.
// If the log level is invalid (or none is specified), no no-op logger will be returned.
func NewLogger() (*zap.Logger, error) {
	switch env.GetLogLevel() {
	case env.LogLevelProd:
		return zap.NewProduction()
	case env.LogLevelDebug:
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}
