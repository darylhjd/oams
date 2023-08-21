package database

import (
	"context"
	"errors"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

func (d *DB) ListUsers(ctx context.Context, params listParams) ([]model.User, error) {
	var res []model.User

	stmt := SELECT(
		Users.AllColumns,
	).FROM(
		Users.Table,
	).ORDER_BY(
		Users.ID.ASC(),
	)

	stmt = setSorts(stmt, params)
	stmt = setLimit(stmt, params)
	stmt = setOffset(stmt, params)

	err := stmt.QueryContext(ctx, d.queryable, &res)
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

	err := stmt.QueryContext(ctx, d.queryable, &res)
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

	stmt := Users.INSERT(
		Users.ID,
		Users.Name,
		Users.Email,
		Users.Role,
	).MODEL(
		model.User{
			ID:    arg.ID,
			Name:  arg.Name,
			Email: arg.Email,
			Role:  arg.Role,
		},
	).RETURNING(
		Users.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

type UpdateUserParams struct {
	Name  *string         `json:"name"`
	Email *string         `json:"email"`
	Role  *model.UserRole `json:"role"`
}

func (d *DB) UpdateUser(ctx context.Context, id string, arg UpdateUserParams) (model.User, error) {
	var (
		res    model.User
		cols   ColumnList
		update model.User
	)

	if arg.Name != nil {
		cols = append(cols, Users.Name)
		update.Name = *arg.Name
	}

	if arg.Email != nil {
		cols = append(cols, Users.Email)
		update.Email = *arg.Email
	}

	if arg.Role != nil {
		cols = append(cols, Users.Role)
		update.Role = *arg.Role
	}

	if len(cols) == 0 {
		return d.GetUser(ctx, id)
	}

	stmt := Users.UPDATE(
		cols,
	).MODEL(
		update,
	).WHERE(
		Users.ID.EQ(String(id)),
	).RETURNING(
		Users.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

func (d *DB) DeleteUser(ctx context.Context, id string) (model.User, error) {
	var res model.User

	stmt := Users.DELETE().
		WHERE(
			Users.ID.EQ(String(id)),
		).
		RETURNING(Users.AllColumns)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

type UpsertUserParams struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// BatchUpsertUsers inserts a batch of users into the database. If the user already exists, then update the name and
// email. Note that the role field is ignored in this operation.
func (d *DB) BatchUpsertUsers(ctx context.Context, args []UpsertUserParams) ([]model.User, error) {
	if len(args) == 0 {
		return nil, nil
	}

	inserts := make([]model.User, 0, len(args))
	{
		dupFinder := map[UpsertUserParams]struct{}{}
		for _, param := range args {
			if _, ok := dupFinder[param]; !ok {
				dupFinder[param] = struct{}{}
				inserts = append(inserts, model.User{
					ID:    param.ID,
					Name:  param.Name,
					Email: param.Email,
				})
			}
		}
	}

	stmt := Users.INSERT(
		Users.ID,
		Users.Name,
		Users.Email,
	).MODELS(
		inserts,
	).ON_CONFLICT(
		Users.ID,
	).DO_UPDATE(
		SET(
			Users.Name.SET(Users.EXCLUDED.Name),
			Users.Email.SET(Users.EXCLUDED.Email),
		),
	).RETURNING(
		Users.AllColumns,
	)

	res := make([]model.User, 0, len(inserts))
	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

type RegisterUserParams struct {
	ID    string
	Email string
}

// RegisterUser is used to register a user into the database on login. If the user already exists, then nothing is done.
func (d *DB) RegisterUser(ctx context.Context, arg RegisterUserParams) (model.User, error) {
	var res model.User

	stmt := Users.INSERT(
		Users.ID,
		Users.Name,
		Users.Email,
	).MODEL(
		model.User{
			ID:    arg.ID,
			Email: arg.Email,
		},
	).ON_CONFLICT(
		Users.ID,
	).DO_NOTHING().RETURNING(Users.AllColumns)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	if err != nil && errors.Is(err, qrm.ErrNoRows) {
		return res, nil
	}

	return res, err
}

type UpcomingClassGroupSession struct {
	Code      string          `alias:"class.code" json:"code"`
	Year      int32           `alias:"class.year" json:"year"`
	Semester  string          `alias:"class.semester" json:"semester"`
	Name      string          `alias:"class_group.name" json:"name"`
	ClassType model.ClassType `alias:"class_group.class_type" json:"class_type"`
	StartTime time.Time       `alias:"class_group_session.start_time" json:"start_time"`
	EndTime   time.Time       `alias:"class_group_session.end_time" json:"end_time"`
	Venue     string          `alias:"class_group_session.venue" json:"venue"`
}

// GetUserUpcomingClassGroupSessions gets information on a user's upcoming classes. This query returns all session
// enrollments for that user that are currently happening or will happen in the future. The sessions are returned in
// ascending order of start time and then end time.
func (d *DB) GetUserUpcomingClassGroupSessions(ctx context.Context, id string) ([]UpcomingClassGroupSession, error) {
	var res []UpcomingClassGroupSession

	stmt := SELECT(
		Classes.Code, Classes.Year, Classes.Semester,
		ClassGroups.Name, ClassGroups.ClassType,
		ClassGroupSessions.StartTime, ClassGroupSessions.EndTime, ClassGroupSessions.Venue,
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

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}
