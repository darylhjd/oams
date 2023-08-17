package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/database"
)

// limitOffsetter allows us to test for migration success.
type limitOffsetter struct{}

func (l limitOffsetter) Limit() *int64 {
	return nil
}

func (l limitOffsetter) Offset() *int64 {
	return nil
}

func TestSetUpTearDown(t *testing.T) {
	ctx := context.Background()
	a := assert.New(t)

	id := uuid.NewString()

	// Run setup.
	testDb := SetUp(t, id)

	// Check that the database is created.
	a.Nil(testDb.Conn.PingContext(ctx))

	// Check that the migration ran correctly.
	_, err := testDb.ListUsers(context.Background(), limitOffsetter{})
	a.Nil(err)

	// Run teardown.
	TearDown(t, testDb, id)

	// Check that the database no longer exists.
	_, err = database.ConnectDB(ctx, id)
	a.ErrorContains(err, "does not exist")
}
