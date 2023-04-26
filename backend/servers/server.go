package servers

import (
	"go.uber.org/zap"
)

// Server defines a basic Server structure for an Oats service.
type Server interface {
	GetLogger() *zap.Logger
}
