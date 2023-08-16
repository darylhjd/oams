package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListUsers(ctx context.Context) ([]model.User, error) {
	var res []model.User

	stmt := SELECT(
		Users.AllColumns,
	).FROM(
		Users.Table,
	).ORDER_BY(
		Users.ID.ASC(),
	)

	err := stmt.QueryContext(ctx, d.Conn, &res)
	return res, err
}
