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

func (d *DB) GetUpcomingManagedClassGroupSessions(ctx context.Context) ([]model.ClassGroupSession, error) {
	var res []model.ClassGroupSession

	stmt := SELECT(
		ClassGroupSessions.AllColumns,
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
