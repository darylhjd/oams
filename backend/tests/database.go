package tests

import (
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/lib/pq"

	"github.com/darylhjd/oats/backend/database"
	_ "github.com/darylhjd/oats/backend/env/autoloader"
)

// TestEnv provides the caller with a sand-boxed test environment.
type TestEnv struct {
	Db  *database.DB
	mig *migrate.Migrate
}

// SetUp is a wrapper function to help set up a new test package.
// This function is useful to create a sandbox database for testing a package.
func SetUp(m *testing.M, namespace string) (*TestEnv, error) {
	// Setup.
	// Create the test database.
	testDb, err := CreateAndConnect(m, namespace)
	if err != nil {
		log.Fatal(err)
	}

	mig, err := database.NewMigrate(namespace, testDb.Db)
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

// RunAndTearDown is a helper function to run and then tear down the given test environment.
func RunAndTearDown(m *testing.M, testEnv *TestEnv, namespace string) {
	res := m.Run()

	sErr, dErr := testEnv.mig.Close()
	if sErr != nil || dErr != nil {
		log.Fatal(sErr, dErr)
	}

	if err := testEnv.Db.Close(); err != nil {
		log.Fatal(err)
	}

	if err := DropDatabase(m, namespace); err != nil {
		log.Fatal(err)
	}

	os.Exit(res)
}

// CreateAndConnect creates a new database with the specified name.
func CreateAndConnect(_ *testing.M, dbName string) (*database.DB, error) {
	db, err := database.ConnectDB("")
	if err != nil {
		return nil, err
	}

	_, err = db.Db.Exec("CREATE DATABASE" + pq.QuoteIdentifier(dbName))
	if err != nil {
		return nil, err
	}

	if err = db.Close(); err != nil {
		return nil, err
	}

	return database.ConnectDB(dbName)
}

// DropDatabase drops a database of the specified name.
func DropDatabase(_ *testing.M, dbName string) error {
	db, err := database.ConnectDB("")
	if err != nil {
		return err
	}

	_, err = db.Db.Exec("DROP DATABASE" + pq.QuoteIdentifier(dbName))
	if err != nil {
		return err
	}

	return db.Close()
}
