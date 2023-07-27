package tests

import (
	"context"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
)

// StubAuthContextUser inserts the mock auth context user into the database.
func StubAuthContextUser(t *testing.T, ctx context.Context, q *database.Queries) {
	t.Helper()

	err := q.UpsertUsers(ctx, []database.UpsertUsersParams{{
		ID:    MockAuthenticatorIDTokenName,
		Email: MockAuthenticatorAccountPreferredUsername,
		Role:  database.UserRoleSTUDENT,
	}}).Close()
	if err != nil {
		t.Fatal(err)
	}
}

// StubUser inserts a mock user with the given ID into the database.
func StubUser(t *testing.T, ctx context.Context, q *database.Queries, id string) {
	t.Helper()

	err := q.UpsertUsers(ctx, []database.UpsertUsersParams{{
		ID:   id,
		Name: "",
		Role: database.UserRoleSTUDENT,
	}}).Close()
	if err != nil {
		t.Fatal(err)
	}
}
