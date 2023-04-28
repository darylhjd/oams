package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/darylhjd/oats/backend/env/autoloader"
)

const (
	migrationNamespace = "migration"

	migrationDir = "config/migrations"
)

//go:embed config/migrations/*.sql
var migrations embed.FS

// NewMigrate is a helper function to get a new migrator for a database using its name or database.
// Either should be non-nil, if both are provided, the default is the db instance.
func NewMigrate(dbName string, db *sql.DB) (*migrate.Migrate, error) {
	// If both not provided.
	if db == nil && dbName == "" {
		return nil, fmt.Errorf("%s - database name or instance not provided", migrationNamespace)
	}

	migrationSource, err := iofs.New(migrations, migrationDir)
	if err != nil {
		return nil, err
	}

	// If database instance provided, use it.
	if db != nil {
		instance, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, err
		}

		return migrate.NewWithInstance(migrationNamespace, migrationSource, Namespace, instance)
	}

	_, connString, err := GetConnectionProperties(dbName)
	if err != nil {
		return nil, err
	}

	return migrate.NewWithSourceInstance(migrationNamespace, migrationSource, connString)
}
