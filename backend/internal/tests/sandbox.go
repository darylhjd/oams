package tests

import (
	"context"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"

	"github.com/darylhjd/oams/backend/internal/database"
)

// TestEnv provides the caller with a sand-boxed test environment.
type TestEnv struct {
	Db  *database.DB
	mig *migrate.Migrate
}

// SetUp is a helper function to help set up a new test package.
// This function is useful to create a sandbox database for testing a package.
func SetUp(_ *testing.M, namespace string) (*TestEnv, error) {
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

	mig, err := database.NewMigrate(namespace, testDb.C)
	if err != nil {
		log.Fatal(err)
	}

	// Do migrations.
	if err = mig.Up(); err != nil {
		log.Fatal(err)
	}

	return &TestEnv{
		testDb, mig,
	}, nil
}

// TearDown is a helper function to tear down the given test environment.
func TearDown(_ *testing.M, testEnv *TestEnv, namespace string) {
	sErr, dErr := testEnv.mig.Close()
	if sErr != nil || dErr != nil {
		log.Fatal(sErr, dErr)
	}

	if err := database.Drop(context.Background(), namespace, true); err != nil {
		log.Fatal(err)
	}
}
