package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListClasses(ctx context.Context, params ListQueryParams) ([]model.Class, error) {
	var res []model.Class

	stmt := SELECT(
		Classes.AllColumns,
	).FROM(
		Classes,
	)

	stmt = params.setFilters(stmt)
	stmt = params.setSorts(stmt)
	stmt = params.setLimit(stmt)
	stmt = params.setOffset(stmt)

	err := stmt.QueryContext(ctx, d.qe, &res)
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

	err := stmt.QueryContext(ctx, d.qe, &res)
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

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type UpsertClassParams struct {
	Code      string `json:"code"`
	Year      int32  `json:"year"`
	Semester  string `json:"semester"`
	Programme string `json:"programme"`
	Au        int16  `json:"au"`
}

// BatchUpsertClasses inserts a batch of classes into the database. If a class already exists, then the programme and
// au fields are updated. This operation is mainly used by the batch endpoint. Note that this operation helps removes any
// potential duplicate UpsertClassParams provided to it.
func (d *DB) BatchUpsertClasses(ctx context.Context, args []UpsertClassParams) ([]model.Class, error) {
	if len(args) == 0 {
		return nil, nil
	}

	inserts := make([]model.Class, 0, len(args))
	{
		dupFinder := map[UpsertClassParams]struct{}{}
		for _, param := range args {
			if _, ok := dupFinder[param]; !ok {
				dupFinder[param] = struct{}{}
				inserts = append(inserts, model.Class{
					Code:      param.Code,
					Year:      param.Year,
					Semester:  param.Semester,
					Programme: param.Programme,
					Au:        param.Au,
				})
			}
		}
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

	res := make([]model.Class, 0, len(inserts))
	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
