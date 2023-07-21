package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/database"
)

func TestSetUpTearDown(t *testing.T) {
	ctx := context.Background()
	a := assert.New(t)

	// Run setup.
	testDb := SetUp(nil, namespace)

	// Check that the database is created.
	a.Nil(testDb.C.Ping(ctx))

	// Check that the migration ran correctly.
	_, err := testDb.Q.ListStudents(context.Background())
	a.Nil(err)

	// Run teardown.
	TearDown(nil, testDb, namespace)

	// Check that the database no longer exists.
	_, err = database.ConnectDB(ctx, namespace)
	a.ErrorContains(err, "does not exist")
}
