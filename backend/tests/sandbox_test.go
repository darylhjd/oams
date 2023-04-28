package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oats/backend/database"
)

func TestSetUpTearDown(t *testing.T) {
	// Run setup.
	testEnv, err := SetUp(nil, namespace)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the database is created.
	if err = testEnv.Db.Db.Ping(); err != nil {
		t.Fatal(err)
	}

	// Check that the migration ran correctly.
	if _, err = testEnv.Db.Q.ListStudents(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Run teardown.
	TearDown(nil, testEnv, namespace)

	// Check that the database no longer exists.
	_, err = database.ConnectDB(namespace)

	a := assert.New(t)
	a.ErrorContains(err, "does not exist")
}
