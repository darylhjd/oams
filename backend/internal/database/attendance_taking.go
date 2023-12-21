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
	ClassType    model.ClassType     `alias:"class_group.class_type"`
	ManagingRole *model.ManagingRole `alias:"class_group_manager.managing_role"` // For nil values, exposed as system admin.
}

func (d *DB) GetUpcomingManagedClassGroupSessions(ctx context.Context) ([]UpcomingManagedClassGroupSession, error) {
	var res []UpcomingManagedClassGroupSession

	stmt := SELECT(
		ClassGroupSessions.ID,
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
		ClassGroupSessions.Venue,
		ClassGroups.ClassType,
		ClassGroupManagers.ManagingRole,
	).FROM(
		ClassGroupSessions.INNER_JOIN(
			ClassGroups, ClassGroups.ID.EQ(ClassGroupSessions.ClassGroupID),
		).LEFT_JOIN(
			ClassGroupManagers, ClassGroupManagers.ClassGroupID.EQ(ClassGroups.ID),
		),
	).WHERE(
		TimestampzT(time.Now()).BETWEEN(
			ClassGroupSessions.StartTime.SUB(INTERVALd(attendanceTimeBuffer)),
			ClassGroupSessions.EndTime,
		).AND(
			attendanceTakingRLS(ctx),
		),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
