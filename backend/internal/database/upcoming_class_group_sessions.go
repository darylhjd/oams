package database

import (
	"context"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

const (
	attendanceTimeBuffer = time.Minute * 30
)

type UpcomingManagedClassGroupSession struct {
	ID           int64               `alias:"class_group_session.id" json:"id"`
	StartTime    time.Time           `alias:"class_group_session.start_time" json:"start_time"`
	EndTime      time.Time           `alias:"class_group_session.end_time" json:"end_time"`
	Venue        string              `alias:"class_group_session.venue" json:"venue"`
	Code         string              `alias:"class.code" json:"code"`
	Year         int32               `alias:"class.year" json:"year"`
	Semester     string              `alias:"class.semester" json:"semester"`
	Name         string              `alias:"class_group.name" json:"name"`
	ClassType    model.ClassType     `alias:"class_group.class_type" json:"class_type"`
	ManagingRole *model.ManagingRole `alias:"class_group_manager.managing_role" json:"managing_role"` // For nil values, exposed as system admin.
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

type AttendanceEntry struct {
	ID        int64  `alias:"session_enrollment.id" json:"id"`
	SessionID int64  `alias:"session_enrollment.session_id" json:"session_id"`
	UserID    string `alias:"session_enrollment.user_id" json:"user_id"`
	UserName  string `alias:"user.name" json:"user_name"`
	Attended  bool   `alias:"session_enrollment.attended" json:"attended"`
}

func (d *DB) GetUpcomingClassGroupAttendanceEntries(ctx context.Context, id int64) ([]AttendanceEntry, error) {
	var res []AttendanceEntry

	stmt := SELECT(
		SessionEnrollments.ID,
		SessionEnrollments.SessionID,
		SessionEnrollments.UserID,
		Users.Name,
		SessionEnrollments.Attended,
	).FROM(
		SessionEnrollments.INNER_JOIN(
			Users, Users.ID.EQ(SessionEnrollments.UserID),
		),
	).WHERE(
		SessionEnrollments.SessionID.EQ(Int64(id)).AND(
			sessionEnrollmentRLS(ctx),
		),
	).ORDER_BY(
		Users.Name,
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type UpdateAttendanceEntryParams struct {
	ClassGroupSessionID int64
	SessionEnrollmentID int64
	UserID              string
	Attended            bool
	UserSignature       string
}

func (d *DB) UpdateAttendanceEntry(ctx context.Context, arg UpdateAttendanceEntryParams) error {
	if oauth2.GetAuthContext(ctx).User.Role != model.UserRole_SystemAdmin {
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
			SessionEnrollments.ID.EQ(Int64(arg.SessionEnrollmentID)),
		)

		err := signatureStmt.QueryContext(ctx, d.qe, &signature)
		if errors.Is(err, qrm.ErrNoRows) {
			if signature.Signature, err = argon2id.CreateHash(arg.UserID, argon2id.DefaultParams); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		match, err := argon2id.ComparePasswordAndHash(arg.UserSignature, signature.Signature)
		if err != nil {
			return err
		} else if !match {
			return qrm.ErrNoRows
		}
	}

	stmt := SessionEnrollments.UPDATE(
		SessionEnrollments.Attended,
	).MODEL(
		model.SessionEnrollment{
			Attended: arg.Attended,
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
								SessionEnrollments.ID.EQ(Int64(arg.SessionEnrollmentID)),
							),
						),
					),
				).AND(
					isManagedClassGroupSession(ctx),
				),
			),
		).AND(
			SessionEnrollments.ID.EQ(Int64(arg.SessionEnrollmentID)),
		),
	)
	_, err := stmt.ExecContext(ctx, d.qe)
	return err
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
	).DISTINCT().FROM(
		ClassGroupSessions.INNER_JOIN(
			ClassGroups, ClassGroups.ID.EQ(ClassGroupSessions.ClassGroupID),
		).INNER_JOIN(
			Classes, Classes.ID.EQ(ClassGroups.ClassID),
		).LEFT_JOIN( // Left Join to support if user is just system admin.
			ClassGroupManagers, ClassGroupManagers.ClassGroupID.EQ(ClassGroups.ID),
		),
	).ORDER_BY(
		ClassGroupSessions.StartTime.ASC(),
	)
}

func isManagedClassGroupSession(ctx context.Context) BoolExpression {
	return TimestampzT(time.Now()).BETWEEN(
		ClassGroupSessions.StartTime.SUB(INTERVALd(attendanceTimeBuffer)),
		ClassGroupSessions.EndTime,
	).AND(
		classGroupManagerRLS(ctx),
	)
}
