package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListClasses(ctx context.Context) ([]model.Class, error) {
	var res []model.Class

	stmt := SELECT(
		Classes.AllColumns,
	).FROM(
		Classes,
	).ORDER_BY(
		Classes.Code,
		Classes.Year,
		Classes.Semester,
	)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

func (d *DB) GetClass(ctx context.Context, id int64) (model.Class, error) {
	var res model.Class

	stmt := SELECT(
		Classes.AllColumns,
	).FROM(
		Classes,
	).WHERE(
		Classes.ID.EQ(Int64(id)),
	).LIMIT(1)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

type CreateClassParams struct {
	Code      string `json:"code"`
	Year      int32  `json:"year"`
	Semester  string `json:"semester"`
	Programme string `json:"programme"`
	Au        int16  `json:"au"`
}

func (d *DB) CreateClass(ctx context.Context, arg CreateClassParams) (model.Class, error) {
	var res model.Class

	now := time.Now()

	stmt := Classes.INSERT(
		Classes.Code,
		Classes.Year,
		Classes.Semester,
		Classes.Programme,
		Classes.Au,
		Classes.CreatedAt,
		Classes.UpdatedAt,
	).MODEL(
		model.Class{
			Code:      arg.Code,
			Year:      arg.Year,
			Semester:  arg.Semester,
			Programme: arg.Programme,
			Au:        arg.Au,
			CreatedAt: now,
			UpdatedAt: now,
		},
	).RETURNING(
		Classes.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

func (d *DB) DeleteClass(ctx context.Context, id int64) (model.Class, error) {
	var res model.Class

	stmt := Classes.DELETE().
		WHERE(
			Classes.ID.EQ(Int64(id)),
		).RETURNING(Classes.AllColumns)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}
