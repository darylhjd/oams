package tests

import (
	"context"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	MockAuthenticatorAccessToken              = "mock-access-token"
	MockAuthenticatorAccountHomeAccountID     = "mock-home-account-id"
	MockAuthenticatorAccountPreferredUsername = "NTU0001@e.ntu.edu.sg"
	MockAuthenticatorIDTokenName              = "TESTACC001"
)

// NewMockUpsertAuthContextUsersParams creates a new UpsertUsersParams that creates the auth context user in the database.
func NewMockUpsertAuthContextUsersParams() database.UpsertUsersParams {
	return database.UpsertUsersParams{
		ID: MockAuthenticatorIDTokenName,
		Email: pgtype.Text{
			String: MockAuthenticatorAccountPreferredUsername,
			Valid:  true,
		},
		Role: database.UserRoleSTUDENT,
	}
}

// StubAuthContextUser inserts the mock auth context user into the database.
func StubAuthContextUser(t *testing.T, ctx context.Context, q *database.Queries) {
	t.Helper()

	err := q.UpsertUsers(ctx, []database.UpsertUsersParams{NewMockUpsertAuthContextUsersParams()}).Close()
	if err != nil {
		t.Fatal(err)
	}
}

// NewMockAuthContextUser creates a new database.User meant for checking that
func NewMockAuthContextUser(createdAt, updatedAt pgtype.Timestamp) database.User {
	return database.User{
		ID:   MockAuthenticatorIDTokenName,
		Name: "",
		Email: pgtype.Text{
			String: MockAuthenticatorAccountPreferredUsername,
			Valid:  true,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
