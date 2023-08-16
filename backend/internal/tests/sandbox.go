package tests

import (
	"context"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
)

// SetUp is a helper function to help set up a new database for a test.
// Use a unique identifier to avoid database clashing.
func SetUp(t *testing.T, id string) *database.DB {
	t.Helper()
	ctx := context.Background()

	// Setup.
	// Create the test database.
	err := database.Create(ctx, id, true)
	if err != nil {
		t.Fatal(err)
	}

	testDb, err := database.ConnectDB(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	mig, err := database.NewMigrate(id)
	if err != nil {
		t.Fatal(err)
	}

	// Do migrations.
	if err = mig.Up(); err != nil {
		t.Fatal(err)
	}

	sErr, dErr := mig.Close()
	if sErr != nil || dErr != nil {
		t.Fatal(sErr, dErr)
	}

	return testDb
}

// TearDown is a helper function to tear down the given test environment.
func TearDown(t *testing.T, db *database.DB, id string) {
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	if err := database.Drop(context.Background(), id, true); err != nil {
		t.Fatal(err)
	}
}
