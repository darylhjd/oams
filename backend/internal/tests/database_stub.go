package tests

import (
	"context"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

// StubAuthContextUser inserts the mock auth context user into the database.
func StubAuthContextUser(t *testing.T, ctx context.Context, q *database.Queries) {
	t.Helper()

	err := q.UpsertUsers(ctx, []database.UpsertUsersParams{{
		ID: MockAuthenticatorIDTokenName,
		Email: pgtype.Text{
			String: MockAuthenticatorAccountPreferredUsername,
			Valid:  true,
		},
		Role: database.UserRoleSTUDENT,
	}}).Close()
	if err != nil {
		t.Fatal(err)
	}
}
