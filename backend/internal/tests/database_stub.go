package tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

// StubClassGroup inserts a mock class and a corresponding class group into the database.
func StubClassGroup(t *testing.T, ctx context.Context, q *database.Queries, name string, classType database.ClassType) database.CreateClassGroupRow {
	t.Helper()

	class := StubClass(t, ctx, q, uuid.NewString(), rand.Int31(), uuid.NewString())

	group, err := q.CreateClassGroup(ctx, database.CreateClassGroupParams{
		ClassID:   class.ID,
		Name:      name,
		ClassType: classType,
	})
	if err != nil {
		t.Fatal(err)
	}

	return group
}

// StubClassGroupWithClassID creates a mock class group using an existing class ID.
func StubClassGroupWithClassID(t *testing.T, ctx context.Context, q *database.Queries, classId int64, name string, classType database.ClassType) database.CreateClassGroupRow {
	t.Helper()

	group, err := q.CreateClassGroup(ctx, database.CreateClassGroupParams{
		ClassID:   classId,
		Name:      name,
		ClassType: classType,
	})
	if err != nil {
		t.Fatal(err)
	}

	return group
}

// StubClassGroupSession inserts a mock class, class group and corresponding class group session into the database.
func StubClassGroupSession(t *testing.T, ctx context.Context, q *database.Queries, startTime, endTime pgtype.Timestamp, venue string) database.CreateClassGroupSessionRow {
	t.Helper()

	classGroup := StubClassGroup(t, ctx, q, uuid.NewString(), database.ClassTypeLEC)

	session, err := q.CreateClassGroupSession(ctx, database.CreateClassGroupSessionParams{
		ClassGroupID: classGroup.ID,
		StartTime:    startTime,
		EndTime:      endTime,
		Venue:        venue,
	})
	if err != nil {
		t.Fatal(err)
	}

	return session
}

// StubClassGroupSessionWithClassGroupID creates a mock class group session using an existing class group ID.
func StubClassGroupSessionWithClassGroupID(t *testing.T, ctx context.Context, q *database.Queries, classGroupId int64, startTime, endTime pgtype.Timestamp, venue string) database.CreateClassGroupSessionRow {
	t.Helper()

	session, err := q.CreateClassGroupSession(ctx, database.CreateClassGroupSessionParams{
		ClassGroupID: classGroupId,
		StartTime:    startTime,
		EndTime:      endTime,
		Venue:        venue,
	})
	if err != nil {
		t.Fatal(err)
	}

	return session
}

// StubSessionEnrollment inserts a mock class group session, user, and corresponding session enrollment into the database.
func StubSessionEnrollment(t *testing.T, ctx context.Context, q *database.Queries, attended bool) database.CreateSessionEnrollmentRow {
	t.Helper()

	session := StubClassGroupSession(t, ctx, q,
		pgtype.Timestamp{Time: time.UnixMicro(1).UTC(), Valid: true},
		pgtype.Timestamp{Time: time.UnixMicro(2).UTC(), Valid: true},
		"VENUE+66",
	)

	user := StubUser(t, ctx, q,
		uuid.NewString(),
		database.UserRoleSTUDENT,
	)

	enrollment, err := q.CreateSessionEnrollment(ctx, database.CreateSessionEnrollmentParams{
		SessionID: session.ID,
		UserID:    user.ID,
		Attended:  attended,
	})
	if err != nil {
		t.Fatal(err)
	}

	return enrollment
}
