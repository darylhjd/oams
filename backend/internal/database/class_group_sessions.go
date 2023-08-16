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

	err := stmt.QueryContext(ctx, d.Conn, &res)
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

	err := stmt.QueryContext(ctx, d.Conn, &res)
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

	now := time.Now()

	stmt := ClassGroupSessions.INSERT(
		ClassGroupSessions.ClassGroupID,
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
		ClassGroupSessions.Venue,
		ClassGroupSessions.CreatedAt,
		ClassGroupSessions.UpdatedAt,
	).MODEL(
		model.ClassGroupSession{
			ClassGroupID: arg.ClassGroupID,
			StartTime:    arg.StartTime,
			EndTime:      arg.EndTime,
			Venue:        arg.Venue,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	).RETURNING(
		ClassGroupSessions.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

func (d *DB) DeleteClassGroupSession(ctx context.Context, id int64) (model.ClassGroupSession, error) {
	var res model.ClassGroupSession

	stmt := ClassGroupSessions.DELETE().
		WHERE(
			ClassGroupSessions.ID.EQ(Int64(id)),
		).RETURNING(ClassGroupSessions.AllColumns)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}
