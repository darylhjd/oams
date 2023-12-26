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
	HasManagedClassGroups bool `alias:".has_managed_class_groups" json:"has_managed_class_groups"`
	IsCourseCoordinator   bool `alias:".is_course_coordinator" json:"is_course_coordinator"`
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
		).AS("has_managed_class_groups"),
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
		).AS("is_course_coordinator"),
	)
	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
