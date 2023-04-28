package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oats/backend/database"
)

func TestSetUpTearDown(t *testing.T) {
	a := assert.New(t)

	// Run setup.
	testEnv, err := SetUp(nil, namespace)
	a.Nil(err)

	// Check that the database is created.
	a.Nil(testEnv.Db.C.Ping())

	// Check that the migration ran correctly.
	_, err = testEnv.Db.Q.ListStudents(context.Background())
	a.Nil(err)

	// Run teardown.
	TearDown(nil, testEnv, namespace)

	// Check that the database no longer exists.
	_, err = database.ConnectDB(namespace)
	a.ErrorContains(err, "does not exist")
}
