package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/logger"
	"go.uber.org/zap"
)

const (
	namespace = "intervention"
)

func main() {
	ctx := context.Background()

	lg, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("%s - failed to initialise logger: %s", namespace, err)
	}

	db, err := database.Connect(ctx)
	if err != nil {
		lg.Fatal(fmt.Sprintf("%s - could not connect to database", namespace), zap.Error(err))
	}
	defer func() {
		if err = db.Close(); err != nil {
			lg.Fatal(fmt.Sprintf("%s - could not close database connection", namespace), zap.Error(err))
		}
	}()

	lg.Info(fmt.Sprintf("%s - ran service", namespace), zap.Time("time", time.Now()))
}
