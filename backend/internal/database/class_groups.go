package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListClassGroups(ctx context.Context) ([]model.ClassGroup, error) {
	var res []model.ClassGroup

	stmt := SELECT(
		ClassGroups.AllColumns,
	).FROM(
		ClassGroups,
	).ORDER_BY(
		ClassGroups.ClassID,
		ClassGroups.Name,
	)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

func (d *DB) GetClassGroup(ctx context.Context, id int64) (model.ClassGroup, error) {
	var res model.ClassGroup

	stmt := SELECT(
		ClassGroups.AllColumns,
	).FROM(
		ClassGroups,
	).WHERE(
		ClassGroups.ID.EQ(Int64(id)),
	).LIMIT(1)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

type CreateClassGroupParams struct {
	ClassID   int64           `json:"class_id"`
	Name      string          `json:"name"`
	ClassType model.ClassType `json:"class_type"`
}

func (d *DB) CreateClassGroup(ctx context.Context, arg CreateClassGroupParams) (model.ClassGroup, error) {
	var res model.ClassGroup

	stmt := ClassGroups.INSERT(
		ClassGroups.ClassID,
		ClassGroups.Name,
		ClassGroups.ClassType,
	).MODEL(
		model.ClassGroup{
			ClassID:   arg.ClassID,
			Name:      arg.Name,
			ClassType: arg.ClassType,
		},
	).RETURNING(
		ClassGroups.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

type UpdateClassGroupParams struct {
	ClassID   *int64           `json:"class_id"`
	Name      *string          `json:"name"`
	ClassType *model.ClassType `json:"class_type"`
}

func (d *DB) UpdateClassGroup(ctx context.Context, id int64, arg UpdateClassGroupParams) (model.ClassGroup, error) {
	var (
		res    model.ClassGroup
		cols   ColumnList
		update model.ClassGroup
	)

	if arg.ClassID != nil {
		cols = append(cols, ClassGroups.ClassID)
		update.ClassID = *arg.ClassID
	}

	if arg.Name != nil {
		cols = append(cols, ClassGroups.Name)
		update.Name = *arg.Name
	}

	if arg.ClassType != nil {
		cols = append(cols, ClassGroups.ClassType)
		update.ClassType = *arg.ClassType
	}

	if len(cols) == 0 {
		return d.GetClassGroup(ctx, id)
	}

	stmt := ClassGroups.UPDATE(
		cols,
	).MODEL(
		update,
	).WHERE(
		ClassGroups.ID.EQ(Int64(id)),
	).RETURNING(
		ClassGroups.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

func (d *DB) DeleteClassGroup(ctx context.Context, id int64) (model.ClassGroup, error) {
	var res model.ClassGroup

	stmt := ClassGroups.DELETE().
		WHERE(
			ClassGroups.ID.EQ(Int64(id)),
		).RETURNING(ClassGroups.AllColumns)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

type UpsertClassGroupParams struct {
	ClassID   int64           `json:"class_id"`
	Name      string          `json:"name"`
	ClassType model.ClassType `json:"class_type"`
}

func (d *DB) UpsertClassGroups(ctx context.Context, args []UpsertClassGroupParams) ([]model.ClassGroup, error) {
	var res []model.ClassGroup

	inserts := make([]model.ClassGroup, 0, len(args))
	for _, param := range args {
		inserts = append(inserts, model.ClassGroup{
			ClassID:   param.ClassID,
			Name:      param.Name,
			ClassType: param.ClassType,
		})
	}

	stmt := ClassGroups.INSERT(
		ClassGroups.ClassID,
		ClassGroups.Name,
		ClassGroups.ClassType,
	).MODELS(
		inserts,
	).ON_CONFLICT().DO_NOTHING().RETURNING(
		ClassGroups.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}
