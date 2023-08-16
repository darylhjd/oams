package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListUsers(ctx context.Context) ([]model.User, error) {
	var res []model.User

	stmt := SELECT(
		Users.AllColumns,
	).FROM(
		Users.Table,
	).ORDER_BY(
		Users.ID.ASC(),
	)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

func (d *DB) GetUser(ctx context.Context, id string) (model.User, error) {
	var res model.User

	stmt := SELECT(
		Users.AllColumns,
	).FROM(
		Users.Table,
	).WHERE(
		Users.ID.EQ(String(id)),
	).LIMIT(1)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

type CreateUserParams struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Email string         `json:"email"`
	Role  model.UserRole `json:"role"`
}

func (d *DB) CreateUser(ctx context.Context, arg CreateUserParams) (model.User, error) {
	var res model.User

	now := time.Now()

	stmt := Users.INSERT(
		Users.AllColumns,
	).MODEL(
		model.User{
			ID:        arg.ID,
			Name:      arg.Name,
			Email:     arg.Email,
			Role:      arg.Role,
			CreatedAt: now,
			UpdatedAt: now,
		},
	).RETURNING(
		Users.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

func (d *DB) DeleteUser(ctx context.Context, id string) (model.User, error) {
	var res model.User

	stmt := Users.DELETE().
		WHERE(
			Users.ID.EQ(String(id)),
		).
		RETURNING(Users.AllColumns)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

type UpcomingClassGroupSession struct {
	model.Class
	model.ClassGroup
	model.ClassGroupSession
}

// GetUserUpcomingClassGroupSessions gets information on a user's upcoming classes. This query returns all session
// enrollments for that user that are currently happening or will happen in the future. The sessions are returned in
// ascending order of start time and then end time.
func (d *DB) GetUserUpcomingClassGroupSessions(ctx context.Context, id string) ([]UpcomingClassGroupSession, error) {
	var res []UpcomingClassGroupSession

	stmt := SELECT(
		Classes.AllColumns, ClassGroups.AllColumns, ClassGroupSessions.AllColumns,
	).FROM(
		ClassGroupSessions.
			INNER_JOIN(ClassGroups, ClassGroups.ID.EQ(ClassGroupSessions.ClassGroupID)).
			INNER_JOIN(Classes, Classes.ID.EQ(ClassGroups.ClassID)),
	).WHERE(
		ClassGroupSessions.ID.IN(
			SELECT(SessionEnrollments.SessionID).
				FROM(SessionEnrollments).
				WHERE(SessionEnrollments.UserID.EQ(String(id))),
		).AND(
			ClassGroupSessions.EndTime.GT(TimestampzT(time.Now())),
		),
	).ORDER_BY(
		ClassGroupSessions.StartTime.ASC(),
		ClassGroupSessions.EndTime.ASC(),
	)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}
