package main

import (
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	"github.com/darylhjd/oats/backend/database"
	_ "github.com/darylhjd/oats/backend/env/autoloader"
	"github.com/darylhjd/oats/backend/logger"
)

//go:embed *.sql
var migrations embed.FS

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("db:migrate - failed to start logger: %s\n", err)
	}

	_, connString, err := database.GetConnectionProperties()
	if err != nil {
		l.Fatal("db:migrate - failed to connect to database", zap.Error(err))
	}

	migrationSource, err := iofs.New(migrations, "")
	if err != nil {
		l.Fatal("db:migrate - failed to load migration files", zap.Error(err))
	}

	m, err := migrate.NewWithSourceInstance("iofs", migrationSource, connString)
	if err != nil {
		l.Fatal("db:migrate - failed to create migration", zap.Error(err))
	}

	if err = m.Up(); err != nil {
		l.Fatal("db:migrate - failed to up-migrate", zap.Error(err))
	}
}
