package database

import (
	"context"

	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/enum"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	. "github.com/go-jet/jet/v2/postgres"
)

func coordinatingClassRLS(ctx context.Context) BoolExpression {
	authContext := oauth2.GetAuthContext(ctx)

	return Bool(
		authContext.User.Role == model.UserRole_SystemAdmin,
	).OR(
		Classes.ID.IN(
			SELECT(
				Classes.ID,
			).FROM(
				Classes.INNER_JOIN(
					ClassGroups, ClassGroups.ClassID.EQ(Classes.ID),
				).INNER_JOIN(
					ClassGroupManagers, ClassGroupManagers.ClassGroupID.EQ(ClassGroups.ID),
				),
			).WHERE(
				ClassGroupManagers.UserID.EQ(String(authContext.User.ID)).AND(
					ClassGroupManagers.ManagingRole.EQ(ManagingRole.CourseCoordinator),
				),
			),
		),
	)
}

func sessionEnrollmentRLS(ctx context.Context) BoolExpression {
	authContext := oauth2.GetAuthContext(ctx)

	return Bool(
		authContext.User.Role == model.UserRole_SystemAdmin,
	).OR(
		SessionEnrollments.UserID.EQ(String(authContext.User.ID)),
	).OR(
		SessionEnrollments.ID.IN(
			SELECT(
				SessionEnrollments.ID,
			).FROM(
				SessionEnrollments.INNER_JOIN(
					ClassGroupSessions, ClassGroupSessions.ID.EQ(SessionEnrollments.SessionID),
				).INNER_JOIN(
					ClassGroups, ClassGroups.ID.EQ(ClassGroupSessions.ClassGroupID),
				).INNER_JOIN(
					ClassGroupManagers, ClassGroupManagers.ClassGroupID.EQ(ClassGroups.ID),
				),
			).WHERE(
				ClassGroupManagers.UserID.EQ(String(authContext.User.ID)),
			),
		),
	)
}
