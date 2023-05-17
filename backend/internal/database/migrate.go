package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/lib/pq"

	_ "github.com/darylhjd/oams/backend/internal/env/autoloader"
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

// NewMigrate is a helper function to get a new migrator for a database using its name or database.
// Either should be non-nil, if both are provided, the default is the db instance.
func NewMigrate(dbName string, db *sql.DB) (*migrate.Migrate, error) {
	// If both not provided.
	if db == nil && dbName == "" {
		return nil, fmt.Errorf("%s - database name or instance not provided", MigrationNamespace)
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

		return migrate.NewWithInstance(MigrationNamespace, migrationSource, Namespace, instance)
	}

	_, connString := GetConnectionProperties(dbName)

	return migrate.NewWithSourceInstance(MigrationNamespace, migrationSource, connString)
}

// Create creates a new database with the specified name.
// Use truncate to specify if the operation deletes an existing database with the same name and creates a new one.
// Warning, this is a high-risk operation!
func Create(dbName string, truncate bool) error {
	db, err := ConnectDB("")
	if err != nil {
		return err
	}

	if truncate {
		if _, err = db.C.Exec(dropDatabaseIfExists + pq.QuoteIdentifier(dbName)); err != nil {
			return err
		}
	}

	if _, err = db.C.Exec(createDatabase + pq.QuoteIdentifier(dbName)); err != nil {
		return err
	}

	return db.Close()
}

// Drop drops a database of the specified name.
// Use mustExist to specify if the operation fails if the database does not exist.
// Warning, this is a high risk operation!
func Drop(dbName string, mustExist bool) error {
	db, err := ConnectDB("")
	if err != nil {
		return err
	}

	statement := dropDatabaseIfExists
	if mustExist {
		statement = dropDatabase
	}

	if _, err = db.C.Exec(statement + pq.QuoteIdentifier(dbName)); err != nil {
		return err
	}

	return db.Close()
}
