package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListClassGroups(ctx context.Context, params ListQueryParams) ([]model.ClassGroup, error) {
	var res []model.ClassGroup

	stmt := SELECT(
		ClassGroups.AllColumns,
	).FROM(
		ClassGroups,
	)

	stmt = params.setSorts(stmt)
	stmt = params.setLimit(stmt)
	stmt = params.setOffset(stmt)

	err := stmt.QueryContext(ctx, d.qe, &res)
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

	err := stmt.QueryContext(ctx, d.qe, &res)
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

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type UpsertClassGroupParams struct {
	ClassID   int64           `json:"class_id"`
	Name      string          `json:"name"`
	ClassType model.ClassType `json:"class_type"`
}

func (d *DB) BatchUpsertClassGroups(ctx context.Context, args []UpsertClassGroupParams) ([]model.ClassGroup, error) {
	if len(args) == 0 {
		return nil, nil
	}

	inserts := make([]model.ClassGroup, 0, len(args))
	{
		dupFinder := map[UpsertClassGroupParams]struct{}{}
		for _, param := range args {
			if _, ok := dupFinder[param]; !ok {
				dupFinder[param] = struct{}{}
				inserts = append(inserts, model.ClassGroup{
					ClassID:   param.ClassID,
					Name:      param.Name,
					ClassType: param.ClassType,
				})
			}
		}
	}

	stmt := ClassGroups.INSERT(
		ClassGroups.ClassID,
		ClassGroups.Name,
		ClassGroups.ClassType,
	).MODELS(
		inserts,
	).ON_CONFLICT().ON_CONSTRAINT(
		"ux_class_id_name_class_type",
	).DO_UPDATE(
		SET(
			ClassGroups.UpdatedAt.SET(TimestampzT(time.Now())),
		),
	).RETURNING(
		ClassGroups.AllColumns,
	)

	res := make([]model.ClassGroup, 0, len(inserts))
	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
