package database

import (
	"context"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/lib/pq"
)

const MigrationNamespace = "database/migration"

const (
	migrationDir         = "config/migrations"
	createDatabase       = "CREATE DATABASE"
	dropDatabaseIfExists = "DROP DATABASE IF EXISTS"
	dropDatabase         = "DROP DATABASE"
)

//go:embed config/migrations/*.sql
var migrations embed.FS

// NewMigrate is a helper function to get a new migrator for a database using its name.
func NewMigrate(dbName string) (*migrate.Migrate, error) {
	// If both not provided.
	if dbName == "" {
		return nil, fmt.Errorf("%s - database name not provided", MigrationNamespace)
	}

	migrationSource, err := iofs.New(migrations, migrationDir)
	if err != nil {
		return nil, err
	}

	_, connString := GetConnectionProperties(dbName)

	return migrate.NewWithSourceInstance(MigrationNamespace, migrationSource, connString)
}

// Create creates a new database with the specified name.
// Use truncate to specify if the operation deletes an existing database with the same name and creates a new one.
// Warning, this is a high-risk operation!
func Create(ctx context.Context, dbName string, truncate bool) error {
	db, err := ConnectDB(ctx, "")
	if err != nil {
		return err
	}
	defer db.Close()

	if truncate {
		if _, err = db.C.Exec(ctx, dropDatabaseIfExists+pq.QuoteIdentifier(dbName)); err != nil {
			return err
		}
	}

	if _, err = db.C.Exec(ctx, createDatabase+pq.QuoteIdentifier(dbName)); err != nil {
		return err
	}

	return nil
}

// Drop drops a database of the specified name.
// Use mustExist to specify if the operation fails if the database does not exist.
// Warning, this is a high risk operation!
func Drop(ctx context.Context, dbName string, mustExist bool) error {
	db, err := ConnectDB(ctx, "")
	if err != nil {
		return err
	}
	defer db.Close()

	statement := dropDatabaseIfExists
	if mustExist {
		statement = dropDatabase
	}

	if _, err = db.C.Exec(ctx, statement+pq.QuoteIdentifier(dbName)); err != nil {
		return err
	}

	return nil
}
