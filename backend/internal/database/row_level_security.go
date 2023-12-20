package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	. "github.com/go-jet/jet/v2/postgres"
)

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

func classGroupManagerRLS(ctx context.Context) BoolExpression {
	authContext := oauth2.GetAuthContext(ctx)

	return Bool(
		authContext.User.Role == model.UserRole_SystemAdmin,
	).OR(
		ClassGroupManagers.UserID.EQ(String(authContext.User.ID)),
	)
}
