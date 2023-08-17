package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListClassGroupSessions(ctx context.Context) ([]model.ClassGroupSession, error) {
	var res []model.ClassGroupSession

	stmt := SELECT(
		ClassGroupSessions.AllColumns,
	).FROM(
		ClassGroupSessions,
	).ORDER_BY(
		ClassGroupSessions.ClassGroupID,
		ClassGroupSessions.StartTime,
	)

	err := stmt.QueryContext(ctx, d.queryable, &res)
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

	err := stmt.QueryContext(ctx, d.queryable, &res)
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

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

type UpdateClassGroupSessionParams struct {
	ClassGroupID *int64  `json:"class_group_id"`
	StartTime    *int64  `json:"start_time"`
	EndTime      *int64  `json:"end_time"`
	Venue        *string `json:"venue"`
}

func (d *DB) UpdateClassGroupSession(ctx context.Context, id int64, arg UpdateClassGroupSessionParams) (model.ClassGroupSession, error) {
	var (
		res    model.ClassGroupSession
		cols   ColumnList
		update model.ClassGroupSession
	)

	if arg.ClassGroupID != nil {
		cols = append(cols, ClassGroupSessions.ClassGroupID)
		update.ClassGroupID = *arg.ClassGroupID
	}

	if arg.StartTime != nil {
		cols = append(cols, ClassGroupSessions.StartTime)
		update.StartTime = time.UnixMicro(*arg.StartTime)
	}

	if arg.EndTime != nil {
		cols = append(cols, ClassGroupSessions.EndTime)
		update.EndTime = time.UnixMicro(*arg.EndTime)
	}

	if arg.Venue != nil {
		cols = append(cols, ClassGroupSessions.Venue)
		update.Venue = *arg.Venue
	}

	if len(cols) == 0 {
		return d.GetClassGroupSession(ctx, id)
	}

	stmt := ClassGroupSessions.UPDATE(
		cols,
	).MODEL(
		update,
	).WHERE(
		ClassGroupSessions.ID.EQ(Int64(id)),
	).RETURNING(
		ClassGroupSessions.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

func (d *DB) DeleteClassGroupSession(ctx context.Context, id int64) (model.ClassGroupSession, error) {
	var res model.ClassGroupSession

	stmt := ClassGroupSessions.DELETE().
		WHERE(
			ClassGroupSessions.ID.EQ(Int64(id)),
		).RETURNING(ClassGroupSessions.AllColumns)

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}

type UpsertClassGroupSessionParams struct {
	ClassGroupID int64     `json:"class_group_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Venue        string    `json:"venue"`
}

func (d *DB) UpsertClassGroupSessions(ctx context.Context, args []UpsertClassGroupSessionParams) ([]model.ClassGroupSession, error) {
	var res []model.ClassGroupSession

	inserts := make([]model.ClassGroupSession, 0, len(args))
	for _, param := range args {
		inserts = append(inserts, model.ClassGroupSession{
			ClassGroupID: param.ClassGroupID,
			StartTime:    param.StartTime,
			EndTime:      param.EndTime,
			Venue:        param.Venue,
		})
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

	err := stmt.QueryContext(ctx, d.queryable, &res)
	return res, err
}
