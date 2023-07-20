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
	testEnv, err := SetUp(nil, namespace)
	a.Nil(err)

	// Check that the database is created.
	a.Nil(testEnv.Db.C.Ping(ctx))

	// Check that the migration ran correctly.
	_, err = testEnv.Db.Q.ListStudents(context.Background())
	a.Nil(err)

	// Run teardown.
	TearDown(nil, testEnv, namespace)

	// Check that the database no longer exists.
	_, err = database.ConnectDB(ctx, namespace)
	a.ErrorContains(err, "does not exist")
}
