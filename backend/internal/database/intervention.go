package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/internal/rules"
	. "github.com/go-jet/jet/v2/postgres"
)

// RuleInfo includes information on each rule.
type RuleInfo struct {
	model.ClassAttendanceRule
	CreatorName   string `alias:"user.name"`
	CreatorEmail  string `alias:"user.email"`
	ClassCode     string `alias:"class.code"`
	ClassYear     int32  `alias:"class.year"`
	ClassSemester string `alias:"class.semester"`
}

// Intervention gets all rules.Fact of classes which had class group sessions occurring today up to now.
// In addition, all RuleInfo of these classes are also returned.
func (d *DB) Intervention(ctx context.Context) ([]rules.Fact, []RuleInfo, error) {
	var facts []rules.Fact

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month()-2, now.Day(), 0, 0, 0, 0, now.Location())

	classGroupSessionPredicate := ClassGroupSessions.StartTime.GT_EQ(TimestampzT(startOfDay)).AND(
		ClassGroupSessions.EndTime.LT(TimestampzT(now)),
	)

	stmt := SELECT(
		Classes.ID,
		Classes.Code,
		Classes.Year,
		Classes.Semester,
		ClassGroups.ClassType,
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
		ClassGroupSessions.Venue,
		SessionEnrollments.UserID,
		Users.Name,
		Users.Email,
		SessionEnrollments.Attended,
	).FROM(
		Classes.INNER_JOIN(
			ClassGroups, ClassGroups.ClassID.EQ(Classes.ID),
		).INNER_JOIN(
			ClassGroupSessions, ClassGroupSessions.ClassGroupID.EQ(ClassGroups.ID),
		).INNER_JOIN(
			SessionEnrollments, SessionEnrollments.SessionID.EQ(ClassGroupSessions.ID),
		).INNER_JOIN(
			Users, Users.ID.EQ(SessionEnrollments.UserID),
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

	var ruleInfos []RuleInfo

	stmt = SELECT(
		ClassAttendanceRules.AllColumns,
		Users.Name,
		Users.Email,
		Classes.Code,
		Classes.Year,
		Classes.Semester,
	).FROM(
		ClassAttendanceRules.INNER_JOIN(
			Classes, Classes.ID.EQ(ClassAttendanceRules.ClassID),
		).INNER_JOIN(
			ClassGroups, ClassGroups.ClassID.EQ(Classes.ID),
		).INNER_JOIN(
			ClassGroupSessions, ClassGroupSessions.ClassGroupID.EQ(ClassGroups.ID),
		).INNER_JOIN(
			Users, Users.ID.EQ(ClassAttendanceRules.CreatorID),
		),
	).WHERE(
		classGroupSessionPredicate.AND(
			ClassAttendanceRules.Active.IS_TRUE(),
		),
	).ORDER_BY(
		Classes.ID,
	)

	err = stmt.QueryContext(ctx, d.qe, &ruleInfos)
	return facts, ruleInfos, err
}
