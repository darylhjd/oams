package database

import (
	"context"

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
		Classes.ID.EQ(Int64(id)).AND(
			coordinatingClassRLS(ctx),
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
		ClassAttendanceRules.ClassID.EQ(Int64(id)).AND(
			classAttendanceRuleRLS(ctx),
		),
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
	).QUERY(
		SELECT(
			Int64(arg.ClassID),
			String(arg.CreatorID),
			String(arg.Title),
			String(arg.Description),
			String(arg.Rule),
			Json(envString),
		).WHERE(
			EXISTS(
				SELECT(
					Classes.AllColumns,
				).FROM(
					Classes,
				).WHERE(
					Classes.ID.EQ(Int64(arg.ClassID)).AND(
						coordinatingClassRLS(ctx),
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

type CoordinatingClassReportData struct {
	Class       model.Class                   `json:"class"`
	Rules       []model.ClassAttendanceRule   `json:"rules"`
	Managers    []ClassGroupManagerReportData `json:"managers"`
	ClassGroups []ClassGroupReportData        `json:"class_groups"`
}

type ClassGroupManagerReportData struct {
	UserID         string             `alias:"user.id" json:"user_id"`
	UserName       string             `alias:"user.name" json:"user_name"`
	ClassGroupName string             `alias:"class_group.name" json:"class_group_name"`
	ManagingRole   model.ManagingRole `alias:"class_group_manager.managing_role" json:"managing_role"`
}

type ClassGroupReportData struct {
	ClassGroup        model.ClassGroup        `json:"class_group"`
	ClassGroupSession model.ClassGroupSession `json:"class_group_session"`
	SessionEnrollment model.SessionEnrollment `json:"session_enrollment"`
}

func (d *DB) GetCoordinatingClassReportData(ctx context.Context, id int64) (CoordinatingClassReportData, error) {
	var res CoordinatingClassReportData

	classStmt := SELECT(
		Classes.AllColumns,
	).FROM(
		Classes,
	).WHERE(
		Classes.ID.EQ(Int64(id)).AND(
			coordinatingClassRLS(ctx),
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
		Classes.ID.EQ(Int64(id)).AND(
			coordinatingClassRLS(ctx),
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
		Classes.ID.EQ(Int64(id)).AND(
			coordinatingClassRLS(ctx),
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
		coordinatingClassRLS(ctx),
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
