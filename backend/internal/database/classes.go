package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListClasses(ctx context.Context, params listParams) ([]model.Class, error) {
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

	stmt = setSorts(stmt, params)
	stmt = setLimit(stmt, params)
	stmt = setOffset(stmt, params)

	err := stmt.QueryContext(ctx, d.queryable, &res)
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

	err := stmt.QueryContext(ctx, d.queryable, &res)
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

	stmt := Classes.INSERT(
		Classes.Code,
		Classes.Year,
		Classes.Semester,
		Classes.Programme,
		Classes.Au,
	).MODEL(
		model.Class{
			Code:      arg.Code,
			Year:      arg.Year,
			Semester:  arg.Semester,
			Programme: arg.Programme,
			Au:        arg.Au,
		},
	).RETURNING(
		Classes.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

type UpdateClassParams struct {
	Code      *string `json:"code"`
	Year      *int32  `json:"year"`
	Semester  *string `json:"semester"`
	Programme *string `json:"programme"`
	Au        *int16  `json:"au"`
}

func (d *DB) UpdateClass(ctx context.Context, id int64, arg UpdateClassParams) (model.Class, error) {
	var (
		res    model.Class
		cols   ColumnList
		update model.Class
	)

	if arg.Code != nil {
		cols = append(cols, Classes.Code)
		update.Code = *arg.Code
	}

	if arg.Year != nil {
		cols = append(cols, Classes.Year)
		update.Year = *arg.Year
	}

	if arg.Semester != nil {
		cols = append(cols, Classes.Semester)
		update.Semester = *arg.Semester
	}

	if arg.Programme != nil {
		cols = append(cols, Classes.Programme)
		update.Programme = *arg.Programme
	}

	if arg.Au != nil {
		cols = append(cols, Classes.Au)
		update.Au = *arg.Au
	}

	if len(cols) == 0 {
		return d.GetClass(ctx, id)
	}

	stmt := Classes.UPDATE(
		cols,
	).MODEL(
		update,
	).WHERE(
		Classes.ID.EQ(Int64(id)),
	).RETURNING(
		Classes.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

func (d *DB) DeleteClass(ctx context.Context, id int64) (model.Class, error) {
	var res model.Class

	stmt := Classes.DELETE().
		WHERE(
			Classes.ID.EQ(Int64(id)),
		).RETURNING(Classes.AllColumns)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

type UpsertClassParams struct {
	Code      string `json:"code"`
	Year      int32  `json:"year"`
	Semester  string `json:"semester"`
	Programme string `json:"programme"`
	Au        int16  `json:"au"`
}

func (d *DB) UpsertClasses(ctx context.Context, args []UpsertClassParams) ([]model.Class, error) {
	var res []model.Class

	inserts := make([]model.Class, 0, len(args))
	for _, param := range args {
		inserts = append(inserts, model.Class{
			Code:      param.Code,
			Year:      param.Year,
			Semester:  param.Semester,
			Programme: param.Programme,
			Au:        param.Au,
		})
	}

	stmt := Classes.INSERT(
		Classes.Code,
		Classes.Year,
		Classes.Semester,
		Classes.Programme,
		Classes.Au,
	).MODELS(
		inserts,
	).ON_CONFLICT().ON_CONSTRAINT(
		"ux_code_year_semester",
	).DO_UPDATE(
		SET(
			Classes.Programme.SET(Classes.EXCLUDED.Programme),
			Classes.Au.SET(Classes.EXCLUDED.Au),
		),
	).RETURNING(
		Classes.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}
