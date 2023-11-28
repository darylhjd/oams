package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListClassManagers(ctx context.Context, params listParams) ([]model.ClassManager, error) {
	var res []model.ClassManager

	stmt := SELECT(
		ClassManagers.AllColumns,
	).FROM(
		ClassManagers,
	)

	stmt = setSorts(stmt, params)
	stmt = setLimit(stmt, params)
	stmt = setOffset(stmt, params)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) GetClassManager(ctx context.Context, id int64) (model.ClassManager, error) {
	var res model.ClassManager

	stmt := SELECT(
		ClassManagers.AllColumns,
	).FROM(
		ClassManagers,
	).WHERE(
		ClassManagers.ID.EQ(Int64(id)),
	).LIMIT(1)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type CreateClassManagerParams struct {
	UserID       string             `json:"user_id"`
	ClassID      int64              `json:"class_id"`
	ManagingRole model.ManagingRole `json:"managing_role"`
}

func (d *DB) CreateClassManager(ctx context.Context, arg CreateClassManagerParams) (model.ClassManager, error) {
	var res model.ClassManager

	stmt := ClassManagers.INSERT(
		ClassManagers.UserID,
		ClassManagers.ClassID,
		ClassManagers.ManagingRole,
	).MODEL(
		model.ClassManager{
			UserID:       arg.UserID,
			ClassID:      arg.ClassID,
			ManagingRole: arg.ManagingRole,
		},
	).RETURNING(
		ClassManagers.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type UpdateClassManagerParams struct {
	ManagingRole *model.ManagingRole `json:"managing_role"`
}

func (d *DB) UpdateClassManager(ctx context.Context, id int64, arg UpdateClassManagerParams) (model.ClassManager, error) {
	var (
		res    model.ClassManager
		cols   ColumnList
		update model.ClassManager
	)

	if arg.ManagingRole != nil {
		cols = append(cols, ClassManagers.ManagingRole)
		update.ManagingRole = *arg.ManagingRole
	}

	if len(cols) == 0 {
		return d.GetClassManager(ctx, id)
	}

	stmt := ClassManagers.UPDATE(
		cols,
	).MODEL(
		update,
	).WHERE(
		ClassManagers.ID.EQ(Int64(id)),
	).RETURNING(
		ClassManagers.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) DeleteClassManager(ctx context.Context, id int64) (model.ClassManager, error) {
	var res model.ClassManager

	stmt := ClassManagers.DELETE().
		WHERE(
			ClassManagers.ID.EQ(Int64(id)),
		).RETURNING(ClassManagers.AllColumns)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
