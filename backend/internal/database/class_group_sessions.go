package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListClassGroupSessions(ctx context.Context, params ListQueryParams) ([]model.ClassGroupSession, error) {
	var res []model.ClassGroupSession

	stmt := SELECT(
		ClassGroupSessions.AllColumns,
	).FROM(
		ClassGroupSessions,
	)

	stmt = params.setSorts(stmt)
	stmt = params.setLimit(stmt)
	stmt = params.setOffset(stmt)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) GetClassGroupSession(ctx context.Context, id int64) (model.ClassGroupSession, error) {
	var res model.ClassGroupSession

	stmt := SELECT(
		ClassGroupSessions.AllColumns,
	).FROM(
		ClassGroupSessions,
	).WHERE(
		ClassGroupSessions.ID.EQ(Int64(id)),
	).LIMIT(1)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type CreateClassGroupSessionParams struct {
	ClassGroupID int64     `json:"class_group_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Venue        string    `json:"venue"`
}

func (d *DB) CreateClassGroupSession(ctx context.Context, arg CreateClassGroupSessionParams) (model.ClassGroupSession, error) {
	var res model.ClassGroupSession

	stmt := ClassGroupSessions.INSERT(
		ClassGroupSessions.ClassGroupID,
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
		ClassGroupSessions.Venue,
	).MODEL(
		model.ClassGroupSession{
			ClassGroupID: arg.ClassGroupID,
			StartTime:    arg.StartTime,
			EndTime:      arg.EndTime,
			Venue:        arg.Venue,
		},
	).RETURNING(
		ClassGroupSessions.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type UpsertClassGroupSessionParams struct {
	ClassGroupID int64     `json:"class_group_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Venue        string    `json:"venue"`
}

// BatchUpsertClassGroupSessions inserts a batch of class group sessions into the database. If a class group session
// already exists, then the end time and venue is updated. This operation is mainly used by the batch endpoint. Note
// that this operation helps remove any potential duplicate UpsertClassGroupSessionParams provided to it.
func (d *DB) BatchUpsertClassGroupSessions(ctx context.Context, args []UpsertClassGroupSessionParams) ([]model.ClassGroupSession, error) {
	if len(args) == 0 {
		return nil, nil
	}

	inserts := make([]model.ClassGroupSession, 0, len(args))
	{
		dupFinder := map[UpsertClassGroupSessionParams]struct{}{}
		for _, param := range args {
			if _, ok := dupFinder[param]; !ok {
				dupFinder[param] = struct{}{}
				inserts = append(inserts, model.ClassGroupSession{
					ClassGroupID: param.ClassGroupID,
					StartTime:    param.StartTime,
					EndTime:      param.EndTime,
					Venue:        param.Venue,
				})
			}
		}
	}

	stmt := ClassGroupSessions.INSERT(
		ClassGroupSessions.ClassGroupID,
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
		ClassGroupSessions.Venue,
	).MODELS(
		inserts,
	).ON_CONFLICT().ON_CONSTRAINT(
		"ux_class_group_id_start_time",
	).DO_UPDATE(
		SET(
			ClassGroupSessions.EndTime.SET(ClassGroupSessions.EXCLUDED.EndTime),
			ClassGroupSessions.Venue.SET(ClassGroupSessions.EXCLUDED.Venue),
		),
	).RETURNING(
		ClassGroupSessions.AllColumns,
	)

	res := make([]model.ClassGroupSession, 0, len(inserts))
	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
