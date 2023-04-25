package server

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/darylhjd/oats/backend/logger"
)

// server defines a basic server structure for an Oats service.
type server struct {
	db *sql.DB
	L  *zap.Logger
}

func newBaseServer() (*server, error) {
	l, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	return &server{
		nil, l,
	}, nil
}
