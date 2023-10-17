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
