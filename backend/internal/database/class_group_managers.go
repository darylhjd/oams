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
	var res []model.ClassGroupManager

	if len(args) == 0 {
		return res, nil
	}

	rowConverter := func(arg UpsertClassGroupManagerParams) SelectStatement {
		return SELECT(
			String(arg.UserID).AS("user_id"),
			String(arg.ClassCode).AS("class_code"),
			Int32(arg.ClassYear).AS("class_year"),
			String(arg.ClassSemester).AS("class_semester"),
			String(arg.ClassGroupName).AS("class_group_name"),
			String(string(arg.ClassType)).AS("class_type"),
			String(string(arg.ManagingRole)).AS("managing_role"),
		)
	}

	// Convert the first row first.
	first := rowConverter(args[0])
	u := UNION(first, first)

	// Convert the reset of the rows.
	for index := 1; index < len(args); index++ {
		u = UNION(u, rowConverter(args[index]))
	}

	tempTable := CTE("arguments")
	userIdCol := StringColumn("user_id").From(tempTable)
	classCodeCol := StringColumn("class_code").From(tempTable)
	classYearCol := IntegerColumn("class_year").From(tempTable)
	classSemCol := StringColumn("class_semester").From(tempTable)
	classGroupNameCol := StringColumn("class_group_name").From(tempTable)
	classTypeCol := StringColumn("class_type").From(tempTable)
	managingRoleCol := StringColumn("managing_role").From(tempTable)

	stmt := WITH(
		tempTable,
	)(
		ClassGroupManagers.INSERT(
			ClassGroupManagers.UserID,
			ClassGroupManagers.ClassGroupID,
			ClassGroupManagers.ManagingRole,
		).QUERY(
			SELECT(
				userIdCol, ClassGroups.ID, managingRoleCol,
			).FROM(
				tempTable.INNER_JOIN(
					Classes, Classes.Code.EQ(classCodeCol).AND(
						Classes.Year.EQ(classYearCol).AND(
							Classes.Semester.EQ(classSemCol),
						),
					),
				).INNER_JOIN(
					ClassGroups, ClassGroups.ClassID.EQ(Classes.ID).AND(
						ClassGroups.Name.EQ(classGroupNameCol).AND(
							ClassGroups.ClassType.EQ(classTypeCol),
						),
					),
				),
			),
		).ON_CONFLICT().ON_CONSTRAINT(
			"ux_user_id_class_group_id",
		).DO_UPDATE(
			SET(
				ClassGroupManagers.ManagingRole.SET(ClassGroupManagers.EXCLUDED.ManagingRole),
			),
		).RETURNING(
			ClassGroupManagers.AllColumns,
		),
	)
	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
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
