package database

import (
	"context"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

const (
	attendanceTimeBuffer = time.Minute * 15
)

type UpcomingManagedClassGroupSession struct {
	ID           int64               `alias:"class_group_session.id"`
	StartTime    time.Time           `alias:"class_group_session.start_time"`
	EndTime      time.Time           `alias:"class_group_session.end_time"`
	Venue        string              `alias:"class_group_session.venue"`
	Code         string              `alias:"class.code"`
	Year         int32               `alias:"class.year"`
	Semester     string              `alias:"class.semester"`
	Name         string              `alias:"class_group.name"`
	ClassType    model.ClassType     `alias:"class_group.class_type"`
	ManagingRole *model.ManagingRole `alias:"class_group_manager.managing_role"` // For nil values, exposed as system admin.
}

func (d *DB) GetUpcomingManagedClassGroupSessions(ctx context.Context) ([]UpcomingManagedClassGroupSession, error) {
	var res []UpcomingManagedClassGroupSession

	stmt := selectManagedClassGroupSessionFields().WHERE(
		isManagedClassGroupSession(ctx),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) GetUpcomingManagedClassGroupSession(ctx context.Context, id int64) (UpcomingManagedClassGroupSession, error) {
	var res UpcomingManagedClassGroupSession

	stmt := selectManagedClassGroupSessionFields().WHERE(
		ClassGroupSessions.ID.EQ(Int64(id)).AND(
			isManagedClassGroupSession(ctx),
		),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) GetUpcomingClassGroupSessionEnrollments(ctx context.Context, id int64) ([]model.SessionEnrollment, error) {
	var res []model.SessionEnrollment

	stmt := SELECT(
		SessionEnrollments.AllColumns,
	).FROM(
		SessionEnrollments,
	).WHERE(
		SessionEnrollments.SessionID.EQ(Int64(id)).AND(
			sessionEnrollmentRLS(ctx),
		),
	).ORDER_BY(
		SessionEnrollments.UserID,
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type UpdateSessionEnrollmentAttendanceParams struct {
	SessionEnrollment   model.SessionEnrollment
	ClassGroupSessionID int64
	UserSignature       string
}

func (d *DB) UpdateSessionEnrollmentAttendance(ctx context.Context, arg UpdateSessionEnrollmentAttendanceParams) (model.SessionEnrollment, error) {
	var res model.SessionEnrollment

	var signature struct {
		Signature string `alias:"user_signature.signature"`
	}
	signatureStmt := SELECT(
		UserSignatures.Signature,
	).FROM(
		UserSignatures.INNER_JOIN(
			SessionEnrollments, SessionEnrollments.UserID.EQ(UserSignatures.UserID),
		),
	).WHERE(
		SessionEnrollments.ID.EQ(Int64(arg.SessionEnrollment.ID)),
	)
	err := signatureStmt.QueryContext(ctx, d.qe, &signature)
	if errors.Is(err, qrm.ErrNoRows) {
		signature.Signature, err = argon2id.CreateHash(arg.SessionEnrollment.UserID, argon2id.DefaultParams)
		if err != nil {
			return res, err
		}
	} else if err != nil {
		return res, err
	}

	match, err := argon2id.ComparePasswordAndHash(arg.UserSignature, signature.Signature)
	if err != nil {
		return res, err
	} else if !match {
		return res, qrm.ErrNoRows
	}

	stmt := SessionEnrollments.UPDATE(
		SessionEnrollments.Attended,
	).MODEL(
		model.SessionEnrollment{
			Attended: arg.SessionEnrollment.Attended,
		},
	).WHERE(
		EXISTS(
			selectManagedClassGroupSessionFields().WHERE(
				ClassGroupSessions.ID.EQ(
					Int64(arg.ClassGroupSessionID),
				).AND(
					ClassGroupSessions.ID.EQ(
						IntExp(
							SELECT(
								SessionEnrollments.SessionID,
							).FROM(
								SessionEnrollments,
							).WHERE(
								SessionEnrollments.ID.EQ(Int64(arg.SessionEnrollment.ID)),
							),
						),
					),
				).AND(
					isManagedClassGroupSession(ctx),
				),
			),
		).AND(
			SessionEnrollments.ID.EQ(Int64(arg.SessionEnrollment.ID)),
		),
	).RETURNING(
		SessionEnrollments.AllColumns,
	)
	err = stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func selectManagedClassGroupSessionFields() SelectStatement {
	return SELECT(
		ClassGroupSessions.ID,
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
		ClassGroupSessions.Venue,
		Classes.Code,
		Classes.Year,
		Classes.Semester,
		ClassGroups.Name,
		ClassGroups.ClassType,
		ClassGroupManagers.ManagingRole,
	).FROM(
		ClassGroupSessions.INNER_JOIN(
			ClassGroups, ClassGroups.ID.EQ(ClassGroupSessions.ClassGroupID),
		).INNER_JOIN(
			Classes, Classes.ID.EQ(ClassGroups.ClassID),
		).LEFT_JOIN(
			ClassGroupManagers, ClassGroupManagers.ClassGroupID.EQ(ClassGroups.ID),
		),
	).ORDER_BY(
		ClassGroupSessions.StartTime.ASC(),
	)
}

func isManagedClassGroupSession(ctx context.Context) BoolExpression {
	// TODO: Uncomment in production.
	//return TimestampzT(time.Now()).BETWEEN(
	//	ClassGroupSessions.StartTime.SUB(INTERVALd(attendanceTimeBuffer)),
	//	ClassGroupSessions.EndTime,
	//).AND(
	return attendanceTakingRLS(ctx)
	//)
}
