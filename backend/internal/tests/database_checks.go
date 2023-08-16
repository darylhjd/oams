package tests

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/stretchr/testify/assert"
)

// CheckUserExists checks that a user with the specified id exists in the test database.
func CheckUserExists(a *assert.Assertions, ctx context.Context, db *database.DB, id string) {
	_, err := db.GetUser(ctx, id)
	a.Nil(err)
}
