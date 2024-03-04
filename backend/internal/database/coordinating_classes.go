package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/internal/rules"
	. "github.com/go-jet/jet/v2/postgres"
)

type CoordinatingClass struct {
	ID        int64  `alias:"class.id" json:"id"`
	Code      string `alias:"class.code" json:"code"`
	Year      int32  `alias:"class.year" json:"year"`
	Semester  string `alias:"class.semester" json:"semester"`
	Programme string `alias:"class.programme" json:"programme"`
	Au        int16  `alias:"class.au" json:"au"`
}

func (d *DB) GetCoordinatingClasses(ctx context.Context) ([]CoordinatingClass, error) {
	var res []CoordinatingClass

	stmt := selectCoordinatingClassFields().WHERE(
		coordinatingClassRLS(ctx),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) GetCoordinatingClass(ctx context.Context, id int64) (CoordinatingClass, error) {
	var res CoordinatingClass

	stmt := selectCoordinatingClassFields().WHERE(
		coordinatingClassRLS(ctx).AND(
			Classes.ID.EQ(Int64(id)),
		),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) GetCoordinatingClassRules(ctx context.Context, id int64) ([]model.ClassAttendanceRule, error) {
	var res []model.ClassAttendanceRule

	stmt := SELECT(
		ClassAttendanceRules.AllColumns,
	).FROM(
		ClassAttendanceRules,
	).WHERE(
		ClassAttendanceRules.ClassID.EQ(Int64(id)),
	).ORDER_BY(
		ClassAttendanceRules.CreatedAt,
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type CreateNewCoordinatingClassRuleParams struct {
	ClassID     int64
	CreatorID   string
	Title       string
	Description string
	Rule        string
	Env         rules.E
}

func (d *DB) CreateNewCoordinatingClassRule(ctx context.Context, arg CreateNewCoordinatingClassRuleParams) (model.ClassAttendanceRule, error) {
	var res model.ClassAttendanceRule

	envString, err := (&rules.Environment{Env: arg.Env}).Value()
	if err != nil {
		return res, err
	}

	stmt := ClassAttendanceRules.INSERT(
		ClassAttendanceRules.ClassID,
		ClassAttendanceRules.CreatorID,
		ClassAttendanceRules.Title,
		ClassAttendanceRules.Description,
		ClassAttendanceRules.Rule,
		ClassAttendanceRules.Environment,
		ClassAttendanceRules.Active,
	).QUERY(
		SELECT(
			Int64(arg.ClassID),
			String(arg.CreatorID),
			String(arg.Title),
			String(arg.Description),
			String(arg.Rule),
			Json(envString),
			Bool(true),
		).WHERE(
			EXISTS(
				SELECT(
					Classes.AllColumns,
				).FROM(
					Classes,
				).WHERE(
					coordinatingClassRLS(ctx).AND(
						Classes.ID.EQ(Int64(arg.ClassID)),
					),
				),
			),
		),
	).RETURNING(
		ClassAttendanceRules.AllColumns,
	)

	err = stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) UpdateCoordinatingClassRule(ctx context.Context, classId, ruleId int64, active bool) (bool, error) {
	var res struct {
		Active bool `alias:"class_attendance_rule.active"`
	}

	stmt := ClassAttendanceRules.UPDATE(
		ClassAttendanceRules.Active,
	).MODEL(
		model.ClassAttendanceRule{
			Active: active,
		},
	).WHERE(
		EXISTS(
			SELECT(
				Classes.AllColumns,
			).FROM(
				Classes,
			).WHERE(
				coordinatingClassRLS(ctx).AND(
					Classes.ID.EQ(Int64(classId)),
				),
			),
		).AND(
			ClassAttendanceRules.ID.EQ(Int64(ruleId)),
		),
	).RETURNING(
		ClassAttendanceRules.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res.Active, err
}

func (d *DB) DeleteCoordinatingClassRule(ctx context.Context, classId, ruleId int64) error {
	var res model.ClassAttendanceRule

	stmt := ClassAttendanceRules.DELETE().
		WHERE(
			EXISTS(
				SELECT(
					Classes.AllColumns,
				).FROM(
					Classes,
				).WHERE(
					coordinatingClassRLS(ctx).AND(
						Classes.ID.EQ(Int64(classId)),
					),
				),
			).AND(
				ClassAttendanceRules.ID.EQ(Int64(ruleId)),
			),
		).RETURNING(
		ClassAttendanceRules.AllColumns,
	)

	return stmt.QueryContext(ctx, d.qe, &res)
}

type CoordinatingClassReportData struct {
	Class       model.Class
	Rules       []model.ClassAttendanceRule
	Managers    []ClassGroupManagerReportData
	ClassGroups []ClassGroupReportData
}

type ClassGroupManagerReportData struct {
	UserID         string             `alias:"user.id"`
	UserName       string             `alias:"user.name"`
	ClassGroupName string             `alias:"class_group.name"`
	ManagingRole   model.ManagingRole `alias:"class_group_manager.managing_role"`
}

type ClassGroupReportData struct {
	ClassGroup        model.ClassGroup
	ClassGroupSession model.ClassGroupSession
	SessionEnrollment model.SessionEnrollment
}

func (d *DB) GetCoordinatingClassReportData(ctx context.Context, id int64) (CoordinatingClassReportData, error) {
	var res CoordinatingClassReportData

	classStmt := SELECT(
		Classes.AllColumns,
	).FROM(
		Classes,
	).WHERE(
		coordinatingClassRLS(ctx).AND(
			Classes.ID.EQ(Int64(id)),
		),
	)
	if err := classStmt.QueryContext(ctx, d.qe, &res.Class); err != nil {
		return res, err
	}

	rulesStmt := SELECT(
		ClassAttendanceRules.AllColumns,
	).FROM(
		ClassAttendanceRules.INNER_JOIN(
			Classes, Classes.ID.EQ(ClassAttendanceRules.ClassID),
		),
	).WHERE(
		coordinatingClassRLS(ctx).AND(
			Classes.ID.EQ(Int64(id)),
		),
	).ORDER_BY(
		ClassAttendanceRules.CreatorID,
	)
	if err := rulesStmt.QueryContext(ctx, d.qe, &res.Rules); err != nil {
		return res, err
	}

	managersStmt := SELECT(
		Users.ID,
		Users.Name,
		ClassGroups.Name,
		ClassGroupManagers.ManagingRole,
	).FROM(
		Classes.INNER_JOIN(
			ClassGroups, ClassGroups.ClassID.EQ(Classes.ID),
		).INNER_JOIN(
			ClassGroupManagers, ClassGroupManagers.ClassGroupID.EQ(ClassGroups.ID),
		).INNER_JOIN(
			Users, Users.ID.EQ(ClassGroupManagers.UserID),
		),
	).WHERE(
		coordinatingClassRLS(ctx).AND(
			Classes.ID.EQ(Int64(id)),
		),
	).ORDER_BY(
		ClassGroups.Name,
		ClassGroupManagers.ManagingRole,
	)
	if err := managersStmt.QueryContext(ctx, d.qe, &res.Managers); err != nil {
		return res, err
	}

	classGroupsStmt := SELECT(
		ClassGroups.AllColumns,
		ClassGroupSessions.AllColumns,
		SessionEnrollments.AllColumns,
	).FROM(
		Classes.INNER_JOIN(
			ClassGroups, ClassGroups.ClassID.EQ(Classes.ID),
		).INNER_JOIN(
			ClassGroupSessions, ClassGroupSessions.ClassGroupID.EQ(ClassGroups.ID),
		).INNER_JOIN(
			SessionEnrollments, SessionEnrollments.SessionID.EQ(ClassGroupSessions.ID),
		),
	).WHERE(
		coordinatingClassRLS(ctx).AND(
			Classes.ID.EQ(Int64(id)),
		).AND(
			ClassGroupSessions.EndTime.LT(TimestampzT(time.Now())),
		),
	).ORDER_BY(
		ClassGroups.Name,
		ClassGroups.ClassType,
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
		SessionEnrollments.UserID,
	)
	if err := classGroupsStmt.QueryContext(ctx, d.qe, &res.ClassGroups); err != nil {
		return res, err
	}

	return res, nil
}

type AttendanceCountData struct {
	ClassGroupName string `alias:"class_group.name" json:"class_group_name"`
	Attended       int    `alias:".attended" json:"attended"`
	NotAttended    int    `alias:".not_attended" json:"not_attended"`
}

func (d *DB) GetDashboardData(ctx context.Context, id int64) ([]AttendanceCountData, error) {
	var res []AttendanceCountData

	stmt := SELECT(
		ClassGroups.Name,
		SUM(CASE().WHEN(SessionEnrollments.Attended.IS_TRUE()).THEN(Int64(1)).ELSE(Int64(0))).AS("attended"),
		SUM(CASE().WHEN(SessionEnrollments.Attended.IS_FALSE()).THEN(Int64(1)).ELSE(Int64(0))).AS("not_attended"),
	).FROM(
		Classes.INNER_JOIN(
			ClassGroups, ClassGroups.ClassID.EQ(Classes.ID),
		).INNER_JOIN(
			ClassGroupSessions, ClassGroupSessions.ClassGroupID.EQ(ClassGroups.ID),
		).INNER_JOIN(
			SessionEnrollments, SessionEnrollments.SessionID.EQ(ClassGroupSessions.ID),
		),
	).WHERE(
		coordinatingClassRLS(ctx).AND(
			Classes.ID.EQ(Int64(id)),
		).AND(
			ClassGroupSessions.EndTime.LT(TimestampzT(time.Now())),
		),
	).GROUP_BY(
		ClassGroups.Name,
	)
	ORDER_BY(
		ClassGroups.Name,
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type ScheduleData struct {
	ClassGroupName      string          `alias:"class_group.name" json:"class_group_name"`
	ClassType           model.ClassType `alias:"class_group.class_type" json:"class_type"`
	ClassGroupSessionID int64           `alias:"class_group_session.id" json:"class_group_session_id"`
	StartTime           time.Time       `alias:"class_group_session.start_time" json:"start_time"`
	EndTime             time.Time       `alias:"class_group_session.end_time" json:"end_time"`
	Venue               string          `alias:"class_group_session.venue" json:"venue"`
}

func (d *DB) GetCoordinatingClassSchedule(ctx context.Context, id int64) ([]ScheduleData, error) {
	var res []ScheduleData

	stmt := SELECT(
		ClassGroups.Name,
		ClassGroups.ClassType,
		ClassGroupSessions.ID,
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
		ClassGroupSessions.Venue,
	).FROM(
		ClassGroupSessions.INNER_JOIN(
			ClassGroups, ClassGroups.ID.EQ(ClassGroupSessions.ClassGroupID),
		).INNER_JOIN(
			Classes, Classes.ID.EQ(ClassGroups.ClassID),
		),
	).WHERE(
		coordinatingClassRLS(ctx).AND(
			Classes.ID.EQ(Int64(id)),
		),
	).ORDER_BY(
		ClassGroupSessions.StartTime,
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

type UpdateCoordinatingClassScheduleParams struct {
	ClassID   int64
	SessionID int64
	StartTime time.Time
	EndTime   time.Time
}

func (d *DB) UpdateCoordinatingClassSchedule(ctx context.Context, arg UpdateCoordinatingClassScheduleParams) (model.ClassGroupSession, error) {
	var res model.ClassGroupSession

	stmt := ClassGroupSessions.UPDATE(
		ClassGroupSessions.StartTime,
		ClassGroupSessions.EndTime,
	).MODEL(
		model.ClassGroupSession{
			StartTime: arg.StartTime,
			EndTime:   arg.EndTime,
		},
	).WHERE(
		ClassGroupSessions.ID.EQ(IntExp(
			SELECT(
				ClassGroupSessions.ID,
			).FROM(
				Classes.INNER_JOIN(
					ClassGroups, ClassGroups.ClassID.EQ(Classes.ID),
				).INNER_JOIN(
					ClassGroupSessions, ClassGroupSessions.ClassGroupID.EQ(ClassGroups.ID),
				),
			).WHERE(
				coordinatingClassRLS(ctx).AND(
					Classes.ID.EQ(Int64(arg.ClassID)),
				).AND(
					ClassGroupSessions.ID.EQ(Int64(arg.SessionID)),
				),
			),
		)),
	).RETURNING(
		ClassGroupSessions.AllColumns,
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func selectCoordinatingClassFields() SelectStatement {
	return SELECT(
		Classes.ID,
		Classes.Code,
		Classes.Year,
		Classes.Semester,
		Classes.Programme,
		Classes.Au,
	).FROM(
		Classes,
	).ORDER_BY(
		Classes.Year.DESC(),
		Classes.Semester.DESC(),
	)
}
