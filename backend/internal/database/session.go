package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/go-jet/jet/v2/postgres"
)

func (d *DB) HasManagedClassGroups(ctx context.Context) (bool, error) {
	var res struct {
		Result bool `alias:"has_managed_class_groups"`
	}

	stmt := postgres.SELECT(
		postgres.EXISTS(
			postgres.SELECT(
				table.ClassGroupManagers.AllColumns,
			).FROM(
				table.ClassGroupManagers,
			).WHERE(
				classGroupManagerRLS(ctx),
			),
		).AS("has_managed_class_groups"),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res.Result, err
}
