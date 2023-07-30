package tests

import (
	"context"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
)

// StubAuthContextUser inserts the mock auth context user into the database.
func StubAuthContextUser(t *testing.T, ctx context.Context, q *database.Queries) {
	t.Helper()

	_, err := q.CreateUser(ctx, database.CreateUserParams{
		ID:    MockAuthenticatorIDTokenName,
		Email: MockAuthenticatorAccountPreferredUsername,
		Role:  database.UserRoleSTUDENT,
	})
	if err != nil {
		t.Fatal(err)
	}
}

// StubUser inserts a mock user with the given ID into the database.
func StubUser(t *testing.T, ctx context.Context, q *database.Queries, id string, role database.UserRole) database.CreateUserRow {
	t.Helper()

	user, err := q.CreateUser(ctx, database.CreateUserParams{
		ID:   id,
		Name: "",
		Role: role,
	})
	if err != nil {
		t.Fatal(err)
	}

	return user
}

// StubClass inserts a mock class with the given fields into the database.
func StubClass(t *testing.T, ctx context.Context, q *database.Queries, code string, year int32, semester string) database.CreateClassRow {
	t.Helper()

	class, err := q.CreateClass(ctx, database.CreateClassParams{
		Code:      code,
		Year:      year,
		Semester:  semester,
		Programme: "",
		Au:        0,
	})
	if err != nil {
		t.Fatal(err)
	}

	return class
}
