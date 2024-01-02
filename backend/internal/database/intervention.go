package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/internal/intervention/fact"
	. "github.com/go-jet/jet/v2/postgres"
)

// Intervention gets all session enrollments of class group sessions that occurred on the current day,
// as well as all the rules for the classes that these class group sessions belong to.
func (d *DB) Intervention(ctx context.Context) ([]fact.F, []model.ClassAttendanceRule, error) {
	var facts []fact.F

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfNextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	classGroupSessionPredicate := ClassGroupSessions.StartTime.GT_EQ(TimestampzT(startOfDay)).AND(
		ClassGroupSessions.EndTime.LT(TimestampzT(startOfNextDay).ADD(INTERVALd(time.Hour * 24 * 20))),
	)

	stmt := SELECT(
		Classes.ID,
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
		ClassGroupSessions.Venue,
		SessionEnrollments.UserID,
		SessionEnrollments.Attended,
	).FROM(
		Classes.INNER_JOIN(
			ClassGroups, ClassGroups.ClassID.EQ(Classes.ID),
		).INNER_JOIN(
			ClassGroupSessions, ClassGroupSessions.ClassGroupID.EQ(ClassGroups.ID),
		).INNER_JOIN(
			SessionEnrollments, SessionEnrollments.SessionID.EQ(ClassGroupSessions.ID),
		),
	).WHERE(
		Classes.ID.IN(
			SELECT(
				Classes.ID,
			).FROM(
				Classes.INNER_JOIN(
					ClassGroups, ClassGroups.ClassID.EQ(Classes.ID),
				).INNER_JOIN(
					ClassGroupSessions, ClassGroupSessions.ClassGroupID.EQ(ClassGroups.ID),
				),
			).WHERE(
				classGroupSessionPredicate,
			),
		),
	).ORDER_BY(
		Classes.ID,
		SessionEnrollments.UserID,
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
	)

	err := stmt.QueryContext(ctx, d.qe, &facts)
	if err != nil {
		return nil, nil, err
	}

	var rules []model.ClassAttendanceRule

	stmt = SELECT(
		ClassAttendanceRules.AllColumns,
	).FROM(
		ClassAttendanceRules.INNER_JOIN(
			Classes, Classes.ID.EQ(ClassAttendanceRules.ClassID),
		).INNER_JOIN(
			ClassGroups, ClassGroups.ClassID.EQ(Classes.ID),
		).INNER_JOIN(
			ClassGroupSessions, ClassGroupSessions.ClassGroupID.EQ(ClassGroups.ID),
		),
	).WHERE(
		classGroupSessionPredicate,
	).ORDER_BY(
		Classes.ID,
	)

	err = stmt.QueryContext(ctx, d.qe, &rules)
	return facts, rules, err
}
