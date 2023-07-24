package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/database"
)

func TestMigrations(t *testing.T) {
	ctx := context.Background()

	// We need to test in order.
	// 1. Create
	// 2. Check creating new migrator
	// 3. Delete
	a := assert.New(t)
	id := uuid.NewString()

	// Make sure this database doesn't currently exist.
	a.Nil(database.Drop(ctx, id, false))

	// Check Create.
	a.Nil(database.Create(ctx, id, false))
	a.Error(database.Create(ctx, id, false))
	a.Nil(database.Create(ctx, id, true))

	// Check Migrator.
	migrator, err := database.NewMigrate(id)
	a.Nil(err)

	sErr, dbErr := migrator.Close()
	if sErr != nil || dbErr != nil {
		t.Fatal(sErr, dbErr)
	}

	// Check if no database name given to NewMigrate.
	_, err = database.NewMigrate("")
	a.ErrorContains(err, "database name not provided")

	// Check dropping.
	a.Nil(database.Drop(ctx, id, true))
	a.Nil(database.Drop(ctx, id, false))
	a.Error(database.Drop(ctx, id, true))
}
