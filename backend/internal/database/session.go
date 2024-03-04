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
	Attendance     bool `alias:".attendance" json:"attendance"`
	Administrative bool `alias:".administrative" json:"administrative"`
}

func (d *DB) GetManagementDetails(ctx context.Context) (ManagementDetails, error) {
	var res ManagementDetails

	auth := oauth2.GetAuthContext(ctx)
	isSystemAdmin := auth.User.Role == model.UserRole_SystemAdmin

	stmt := SELECT(
		EXISTS(
			SELECT(
				ClassGroupManagers.AllColumns,
			).FROM(
				ClassGroupManagers,
			).WHERE(
				ClassGroupManagers.UserID.EQ(String(auth.User.ID)),
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
				ClassGroupManagers.UserID.EQ(String(auth.User.ID)).AND(
					ClassGroupManagers.ManagingRole.EQ(ManagingRole.CourseCoordinator),
				),
			),
		).OR(
			Bool(isSystemAdmin),
		).AS("administrative"),
	)
	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
