package tests

import (
	"context"
	"log"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
)

// SetUp is a helper function to help set up a new test package.
// This function is useful to create a sandbox database for testing a package.
func SetUp(_ *testing.M, namespace string) *database.DB {
	ctx := context.Background()

	// Setup.
	// Create the test database.
	err := database.Create(ctx, namespace, true)
	if err != nil {
		log.Fatal(err)
	}

	testDb, err := database.ConnectDB(ctx, namespace)
	if err != nil {
		log.Fatal(err)
	}

	mig, err := database.NewMigrate(namespace)
	if err != nil {
		log.Fatal(err)
	}

	// Do migrations.
	if err = mig.Up(); err != nil {
		log.Fatal(err)
	}

	sErr, dErr := mig.Close()
	if sErr != nil || dErr != nil {
		log.Fatal(sErr, dErr)
	}

	return testDb
}

// TearDown is a helper function to tear down the given test environment.
func TearDown(_ *testing.M, db *database.DB, namespace string) {
	db.Close()

	if err := database.Drop(context.Background(), namespace, true); err != nil {
		log.Fatal(err)
	}
}
