package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListSessionEnrollments(ctx context.Context, params ListQueryParams) ([]model.SessionEnrollment, error) {
	var res []model.SessionEnrollment

	stmt := SELECT(
		SessionEnrollments.AllColumns,
	).FROM(
		SessionEnrollments,
	).WHERE(
		sessionEnrollmentRLS(ctx),
	)

	stmt = params.setFilters(stmt)
	stmt = params.setSorts(stmt)
	stmt = params.setLimit(stmt)
	stmt = params.setOffset(stmt)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) GetSessionEnrollment(ctx context.Context, id int64) (model.SessionEnrollment, error) {
	var res model.SessionEnrollment

	stmt := SELECT(
		SessionEnrollments.AllColumns,
	).FROM(
		SessionEnrollments,
	).WHERE(
		SessionEnrollments.ID.EQ(Int64(id)).AND(
			sessionEnrollmentRLS(ctx),
		),
	).LIMIT(1)

	err := stmt.QueryContext(ctx, d.qe, &res)
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

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type UpsertSessionEnrollmentParams struct {
	SessionID int64  `json:"session_id"`
	UserID    string `json:"user_id"`
	Attended  bool   `json:"attended"`
}

// BatchUpsertSessionEnrollments inserts a batch of session enrollments into the database. If the session enrollment
// already exists, then nothing is done. Note that the attended field is ignored in this operation.
func (d *DB) BatchUpsertSessionEnrollments(ctx context.Context, args []UpsertSessionEnrollmentParams) ([]model.SessionEnrollment, error) {
	if len(args) == 0 {
		return nil, nil
	}

	inserts := make([]model.SessionEnrollment, 0, len(args))
	{
		dupFinder := map[UpsertSessionEnrollmentParams]struct{}{}
		for _, param := range args {
			if _, ok := dupFinder[param]; !ok {
				dupFinder[param] = struct{}{}
				inserts = append(inserts, model.SessionEnrollment{
					SessionID: param.SessionID,
					UserID:    param.UserID,
					Attended:  param.Attended,
				})
			}
		}
	}

	stmt := SessionEnrollments.INSERT(
		SessionEnrollments.SessionID,
		SessionEnrollments.UserID,
		SessionEnrollments.Attended,
	).MODELS(
		inserts,
	).ON_CONFLICT().ON_CONSTRAINT(
		"ux_session_id_user_id",
	).DO_UPDATE(
		SET(
			SessionEnrollments.UpdatedAt.SET(TimestampzT(time.Now())),
		),
	).RETURNING(
		SessionEnrollments.AllColumns,
	)

	res := make([]model.SessionEnrollment, 0, len(inserts))
	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
