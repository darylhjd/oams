package database

import (
	"context"
	"errors"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

func (d *DB) ListUsers(ctx context.Context, params ListQueryParams) ([]model.User, error) {
	var res []model.User

	stmt := SELECT(
		Users.AllColumns,
	).FROM(
		Users,
	).ORDER_BY(
		Users.ID.ASC(),
	)

	stmt = params.setSorts(stmt)
	stmt = params.setLimit(stmt)
	stmt = params.setOffset(stmt)

	err := stmt.QueryContext(ctx, d.qe, &res)
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

	err := stmt.QueryContext(ctx, d.qe, &res)
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

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type UpsertUserParams struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// BatchUpsertUsers inserts a batch of users into the database. If the user already exists, then update the name and
// email.
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
					ID:   param.ID,
					Name: param.Name,
					Role: model.UserRole_User,
				})
			}
		}
	}

	stmt := Users.INSERT(
		Users.ID,
		Users.Name,
		Users.Email,
		Users.Role,
	).MODELS(
		inserts,
	).ON_CONFLICT(
		Users.ID,
	).DO_UPDATE(
		SET(
			Users.Name.SET(Users.EXCLUDED.Name),
		),
	).RETURNING(
		Users.AllColumns,
	)

	res := make([]model.User, 0, len(inserts))
	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type RegisterUserParams struct {
	ID    string
	Email string
	Role  model.UserRole
}

// RegisterUser is used to register a user into the database on login. If the user already exists, then nothing is done.
func (d *DB) RegisterUser(ctx context.Context, arg RegisterUserParams) (model.User, error) {
	var res model.User

	stmt := Users.INSERT(
		Users.ID,
		Users.Name,
		Users.Email,
		Users.Role,
	).MODEL(
		model.User{
			ID:    arg.ID,
			Email: arg.Email,
			Role:  arg.Role,
		},
	).ON_CONFLICT(
		Users.ID,
	).DO_UPDATE(
		SET(
			Users.Email.SET(Users.EXCLUDED.Email),
		).WHERE(
			Users.Email.NOT_EQ(Users.EXCLUDED.Email),
		),
	).RETURNING(Users.AllColumns)

	err := stmt.QueryContext(ctx, d.qe, &res)
	if err != nil && errors.Is(err, qrm.ErrNoRows) {
		return d.GetUser(ctx, arg.ID)
	}

	return res, err
}
