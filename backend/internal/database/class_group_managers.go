package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListClassGroupManagers(ctx context.Context, params ListQueryParams) ([]model.ClassGroupManager, error) {
	var res []model.ClassGroupManager

	stmt := SELECT(
		ClassGroupManagers.AllColumns,
	).FROM(
		ClassGroupManagers,
	).WHERE(
		classGroupManagerRLS(ctx),
	)

	stmt = params.setFilters(stmt)
	stmt = params.setSorts(stmt)
	stmt = params.setLimit(stmt)
	stmt = params.setOffset(stmt)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type CreateClassGroupManagerParams struct {
	UserID       string             `json:"user_id"`
	ClassGroupID int64              `json:"class_group_id"`
	ManagingRole model.ManagingRole `json:"managing_role"`
}

func (d *DB) CreateClassGroupManager(ctx context.Context, arg CreateClassGroupManagerParams) (model.ClassGroupManager, error) {
	var res model.ClassGroupManager

	stmt := ClassGroupManagers.INSERT(
		ClassGroupManagers.UserID,
		ClassGroupManagers.ClassGroupID,
		ClassGroupManagers.ManagingRole,
	).MODEL(
		model.ClassGroupManager{
			UserID:       arg.UserID,
			ClassGroupID: arg.ClassGroupID,
			ManagingRole: arg.ManagingRole,
		},
	).RETURNING(
		ClassGroupManagers.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type UpsertClassGroupManagerParams struct {
	UserID         string             `json:"user_id"`
	ClassCode      string             `json:"class_code"`
	ClassYear      int32              `json:"class_year"`
	ClassSemester  string             `json:"class_semester"`
	ClassGroupName string             `json:"class_group_name"`
	ClassType      model.ClassType    `json:"class_type"`
	ManagingRole   model.ManagingRole `json:"managing_role"`
}

func (d *DB) BatchUpsertClassGroupManagers(ctx context.Context, args []UpsertClassGroupManagerParams) ([]model.ClassGroupManager, error) {
	// TODO: Implement SQL.
	return nil, nil
}

func (d *DB) HasManagedClassGroups(ctx context.Context) (bool, error) {
	var res struct {
		Result bool `alias:"has_managed_class_groups"`
	}

	stmt := SELECT(
		EXISTS(
			SELECT(
				ClassGroupManagers.AllColumns,
			).FROM(
				ClassGroupManagers,
			).WHERE(
				classGroupManagerRLS(ctx),
			),
		).AS("has_managed_class_groups"),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res.Result, err
}
