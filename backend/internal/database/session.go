package database

import (
	"context"

	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/enum"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	. "github.com/go-jet/jet/v2/postgres"
)

type ManagementDetails struct {
	Attendance      bool `alias:".attendance" json:"attendance"`
	RulesAndReports bool `alias:".rules_and_reports" json:"rules_and_reports"`
}

func (d *DB) GetManagementDetails(ctx context.Context) (ManagementDetails, error) {
	var res ManagementDetails

	isSystemAdmin := oauth2.GetAuthContext(ctx).User.Role == model.UserRole_SystemAdmin

	stmt := SELECT(
		EXISTS(
			SELECT(
				ClassGroupManagers.AllColumns,
			).FROM(
				ClassGroupManagers,
			).WHERE(
				classGroupManagerRLS(ctx),
			),
		).OR(
			Bool(isSystemAdmin),
		).AS("attendance"),
		EXISTS(
			SELECT(
				ClassGroupManagers.AllColumns,
			).FROM(
				ClassGroupManagers,
			).WHERE(
				classGroupManagerRLS(ctx).AND(
					ClassGroupManagers.ManagingRole.EQ(ManagingRole.CourseCoordinator),
				),
			),
		).OR(
			Bool(isSystemAdmin),
		).AS("rules_and_reports"),
	)
	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
