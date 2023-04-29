package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oats/backend/database"
)

const migrationTest = "migration_test"

var migratorTestDb *database.DB

func TestMigrations(t *testing.T) {
	// We need to test in order.
	// 1. Create
	// 2. Check creating new migrator
	// 3. Delete
	a := assert.New(t)

	// Make sure this database doesn't currently exist.
	a.Nil(database.Drop(migrationTest, false))

	// Check Create.
	a.Nil(database.Create(migrationTest, false))
	a.Error(database.Create(migrationTest, false))
	a.Nil(database.Create(migrationTest, true))

	var err error
	migratorTestDb, err = database.ConnectDB(migrationTest)
	a.Nil(err)

	// Check Migrator.
	migrator, err := database.NewMigrate(database.Namespace, nil)
	a.Nil(err)

	sErr, dbErr := migrator.Close()
	if sErr != nil || dbErr != nil {
		t.Fatal(sErr, dbErr)
	}

	// This time, use database instance to get migrator.
	migrator, err = database.NewMigrate("", migratorTestDb.C)
	a.Nil(err)

	sErr, dbErr = migrator.Close()
	if sErr != nil || dbErr != nil {
		t.Fatal(sErr, dbErr)
	}

	// Check if no valid arguments to NewMigrate.
	_, err = database.NewMigrate("", nil)
	a.ErrorContains(err, "database name or instance not provided")

	// Check dropping.
	a.Nil(database.Drop(migrationTest, true))
	a.Nil(database.Drop(migrationTest, false))
	a.Error(database.Drop(migrationTest, true))
}
