package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
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
	ClassType    model.ClassType     `alias:"class_group.class_type"`
	ManagingRole *model.ManagingRole `alias:"class_group_manager.managing_role"` // For nil values, exposed as system admin.
}

func (d *DB) GetUpcomingManagedClassGroupSessions(ctx context.Context) ([]UpcomingManagedClassGroupSession, error) {
	var res []UpcomingManagedClassGroupSession

	stmt := selectManagedClassGroupSession(ctx)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) GetUpcomingManagedClassGroupSession(ctx context.Context, id int64) (UpcomingManagedClassGroupSession, error) {
	var res UpcomingManagedClassGroupSession

	stmt := selectManagedClassGroupSession(ctx).WHERE(
		ClassGroupSessions.ID.EQ(Int64(id)),
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
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func selectManagedClassGroupSession(ctx context.Context) SelectStatement {
	return SELECT(
		ClassGroupSessions.ID,
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
		ClassGroupSessions.Venue,
		Classes.Code,
		Classes.Year,
		Classes.Semester,
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
	).WHERE(
		// TODO: Uncomment in production.
		//TimestampzT(time.Now()).BETWEEN(
		//	ClassGroupSessions.StartTime.SUB(INTERVALd(attendanceTimeBuffer)),
		//	ClassGroupSessions.EndTime,
		//).AND(
		attendanceTakingRLS(ctx),
		//),
	).ORDER_BY(
		ClassGroupSessions.StartTime.ASC(),
	)
}
