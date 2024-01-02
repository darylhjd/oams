package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/pkg/datetime"
	. "github.com/go-jet/jet/v2/postgres"
)

type InterventionData struct {
	ClassGroup        model.ClassGroup
	ClassGroupSession model.ClassGroupSession
	SessionEnrollment model.SessionEnrollment
}

// Intervention gets all session enrollments of class group sessions that occurred in the past.
func (d *DB) Intervention(ctx context.Context) ([]InterventionData, error) {
	var res []InterventionData

	now := time.Now()
	startOfDay := Timestampz(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, datetime.Location.String())
	endOfDay := Timestampz(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, datetime.Location.String())

	stmt := SELECT(
		ClassGroups.AllColumns,
		ClassGroupSessions.AllColumns,
		SessionEnrollments.AllColumns,
	).FROM(
		ClassGroupSessions.INNER_JOIN(
			SessionEnrollments, SessionEnrollments.SessionID.EQ(ClassGroupSessions.ID),
		).INNER_JOIN(
			ClassGroups, ClassGroups.ID.EQ(ClassGroupSessions.ClassGroupID),
		),
	).WHERE(
		ClassGroupSessions.StartTime.GT_EQ(startOfDay).AND(
			ClassGroupSessions.EndTime.LT(endOfDay),
		),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
