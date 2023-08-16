package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListSessionEnrollments(ctx context.Context) ([]model.SessionEnrollment, error) {
	var res []model.SessionEnrollment

	stmt := SELECT(
		SessionEnrollments.AllColumns,
	).FROM(
		SessionEnrollments,
	).ORDER_BY(
		SessionEnrollments.SessionID,
		SessionEnrollments.UserID,
	)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

func (d *DB) GetSessionEnrollment(ctx context.Context, id int64) (model.SessionEnrollment, error) {
	var res model.SessionEnrollment

	stmt := SELECT(
		SessionEnrollments.AllColumns,
	).FROM(
		SessionEnrollments,
	).WHERE(
		SessionEnrollments.ID.EQ(Int64(id)),
	).LIMIT(1)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

type CreateSessionEnrollmentParams struct {
	SessionID int64  `json:"session_id"`
	UserID    string `json:"user_id"`
	Attended  bool   `json:"attended"`
}

func (d *DB) CreateSessionEnrollment(ctx context.Context, arg CreateSessionEnrollmentParams) (model.SessionEnrollment, error) {
	var res model.SessionEnrollment

	stmt := SessionEnrollments.INSERT(
		SessionEnrollments.SessionID,
		SessionEnrollments.UserID,
		SessionEnrollments.Attended,
	).MODEL(
		model.SessionEnrollment{
			SessionID: arg.SessionID,
			UserID:    arg.UserID,
			Attended:  arg.Attended,
		},
	).RETURNING(
		SessionEnrollments.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}

func (d *DB) DeleteSessionEnrollment(ctx context.Context, id int64) (model.SessionEnrollment, error) {
	var res model.SessionEnrollment

	stmt := SessionEnrollments.DELETE().
		WHERE(
			SessionEnrollments.ID.EQ(Int64(id)),
		).RETURNING(SessionEnrollments.AllColumns)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}
