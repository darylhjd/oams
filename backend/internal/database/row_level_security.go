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

// sessionEnrollmentRLS scopes session enrollments entities to be shown only to users with appropriate privileges.
// The following users have access to the relevant scopes:
// 1. System Administrators have access to all records.
// 2. External Systems have access to all records.
// 3. Course Coordinators have access to enrollment data in the context of their coordinating courses.
// 4. Teaching Assistants have access to enrollment data in the context of the courses they are assisting in.
// 5. Normal users have access to their own enrollment data.
func sessionEnrollmentRLS(ctx context.Context) BoolExpression {
	auth := oauth2.GetAuthContext(ctx)

	return Bool(
		auth.User.Role == model.UserRole_SystemAdmin || auth.User.Role == model.UserRole_ExternalService,
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
				ClassGroupManagers.UserID.EQ(String(auth.User.ID)),
			),
		),
	).OR(
		SessionEnrollments.UserID.EQ(String(auth.User.ID)),
	)
}
