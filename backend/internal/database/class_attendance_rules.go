package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) ListClassAttendanceRules(ctx context.Context, params ListQueryParams) ([]model.ClassAttendanceRule, error) {
	var res []model.ClassAttendanceRule

	stmt := SELECT(
		ClassAttendanceRules.AllColumns,
	).FROM(
		ClassAttendanceRules,
	).WHERE(
		classAttendanceRuleRLS(ctx),
	)

	stmt = params.setSorts(stmt)
	stmt = params.setLimit(stmt)
	stmt = params.setOffset(stmt)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
