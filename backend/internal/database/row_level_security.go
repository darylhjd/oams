package database

import (
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/enum"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

// whereIsSystemAdmin adds a boolean check to see if the given user is a system admin. This is useful for enforcing
// row level security.
func whereIsSystemAdmin(userId string) BoolExpression {
	return BoolExp(
		SELECT(
			Users.Role.EQ(UserRole.SystemAdmin),
		).FROM(
			Users,
		).WHERE(
			Users.ID.EQ(String(userId)),
		).LIMIT(1),
	).IS_TRUE()
}
