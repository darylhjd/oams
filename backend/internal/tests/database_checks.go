package tests

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/stretchr/testify/assert"
)

func CheckUserExists(a *assert.Assertions, ctx context.Context, q *database.Queries, id string) {
	_, err := q.GetUser(ctx, id)
	a.Nil(err)
}
