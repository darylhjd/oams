package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	"github.com/darylhjd/oams/backend/internal/middleware/values"
	. "github.com/go-jet/jet/v2/postgres"
)

func sessionEnrollmentRLS(ctx context.Context) BoolExpression {
	authContext := values.GetAuthContext(ctx)

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
					Classes, Classes.ID.EQ(ClassGroups.ClassID),
				).INNER_JOIN(
					ClassManagers, ClassManagers.ClassID.EQ(Classes.ID),
				),
			).WHERE(
				ClassManagers.UserID.EQ(String(authContext.User.ID)),
			),
		),
	)
}

func classManagerRLS(ctx context.Context) BoolExpression {
	authContext := values.GetAuthContext(ctx)

	return Bool(
		authContext.User.Role == model.UserRole_SystemAdmin,
	).OR(
		ClassManagers.UserID.EQ(String(authContext.User.ID)),
	)
}
