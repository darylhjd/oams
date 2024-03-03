package tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/google/uuid"
)

// StubUser inserts a mock user with the given ID into the database.
func StubUser(t *testing.T, ctx context.Context, db *database.DB, params database.CreateUserParams) model.User {
	t.Helper()

	user, err := db.CreateUser(ctx, params)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

// StubClass inserts a mock class with the given fields into the database.
func StubClass(t *testing.T, ctx context.Context, db *database.DB, params database.CreateClassParams) model.Class {
	t.Helper()

	class, err := db.CreateClass(ctx, params)
	if err != nil {
		t.Fatal(err)
	}

	return class
}

// StubClassGroupManager inserts a mock user and class group into the database, and then makes the user
// a manager of the mocked class group.
func StubClassGroupManager(t *testing.T, ctx context.Context, db *database.DB, role model.ManagingRole, classType model.ClassType) model.ClassGroupManager {
	t.Helper()

	user := StubUser(t, ctx, db, database.CreateUserParams{
		ID:   uuid.NewString(),
		Role: model.UserRole_User,
	})
	classGroup := StubClassGroup(t, ctx, db, uuid.NewString(), classType)

	manager, err := db.CreateClassGroupManager(ctx, database.CreateClassGroupManagerParams{
		UserID:       user.ID,
		ClassGroupID: classGroup.ID,
		ManagingRole: role,
	})
	if err != nil {
		t.Fatal(err)
	}

	return manager
}

// StubClassGroup inserts a mock class and a corresponding class group into the database.
func StubClassGroup(t *testing.T, ctx context.Context, db *database.DB, name string, classType model.ClassType) model.ClassGroup {
	t.Helper()

	class := StubClass(t, ctx, db, database.CreateClassParams{
		Code:     uuid.NewString(),
		Year:     rand.Int31(),
		Semester: uuid.NewString(),
	})

	group, err := db.CreateClassGroup(ctx, database.CreateClassGroupParams{
		ClassID:   class.ID,
		Name:      name,
		ClassType: classType,
	})
	if err != nil {
		t.Fatal(err)
	}

	return group
}

// StubClassGroupSession inserts a mock class, class group and corresponding class group session into the database.
func StubClassGroupSession(t *testing.T, ctx context.Context, db *database.DB, startTime, endTime time.Time, venue string) model.ClassGroupSession {
	t.Helper()

	classGroup := StubClassGroup(t, ctx, db, uuid.NewString(), model.ClassType_Lec)

	session, err := db.CreateClassGroupSession(ctx, database.CreateClassGroupSessionParams{
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

// StubSessionEnrollment inserts a mock class group session, user, and corresponding session enrollment into the database.
func StubSessionEnrollment(t *testing.T, ctx context.Context, db *database.DB, attended bool) model.SessionEnrollment {
	t.Helper()

	user := StubUser(t, ctx, db, database.CreateUserParams{
		ID:   uuid.NewString(),
		Role: model.UserRole_User,
	})

	session := StubClassGroupSession(t, ctx, db,
		time.UnixMicro(1),
		time.UnixMicro(2),
		uuid.NewString(),
	)

	enrollment, err := db.CreateSessionEnrollment(ctx, database.CreateSessionEnrollmentParams{
		SessionID: session.ID,
		UserID:    user.ID,
		Attended:  attended,
	})
	if err != nil {
		t.Fatal(err)
	}

	return enrollment
}
