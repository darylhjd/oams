package database

import (
	"context"
	"time"

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

	err := stmt.QueryContext(ctx, d.Conn, &res)
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

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

type CreateClassGroupParams struct {
	ClassID   int64           `json:"class_id"`
	Name      string          `json:"name"`
	ClassType model.ClassType `json:"class_type"`
}

func (d *DB) CreateClassGroup(ctx context.Context, arg CreateClassGroupParams) (model.ClassGroup, error) {
	var res model.ClassGroup

	now := time.Now()

	stmt := ClassGroups.INSERT(
		ClassGroups.ClassID,
		ClassGroups.Name,
		ClassGroups.ClassType,
		ClassGroups.CreatedAt,
		ClassGroups.UpdatedAt,
	).MODEL(
		model.ClassGroup{
			ClassID:   arg.ClassID,
			Name:      arg.Name,
			ClassType: arg.ClassType,
			CreatedAt: now,
			UpdatedAt: now,
		},
	).RETURNING(
		ClassGroups.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

func (d *DB) DeleteClassGroup(ctx context.Context, id int64) (model.ClassGroup, error) {
	var res model.ClassGroup

	stmt := ClassGroups.DELETE().
		WHERE(
			ClassGroups.ID.EQ(Int64(id)),
		).RETURNING(ClassGroups.AllColumns)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}
