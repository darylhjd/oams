package database

import (
	"context"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/pkg/datetime"
	. "github.com/go-jet/jet/v2/postgres"
)

type InterventionData struct {
	model.ClassGroupSession
}

func (d *DB) GetTodayClassGroupSessions(ctx context.Context) ([]InterventionData, error) {
	var res []InterventionData

	now := time.Now()
	startOfDay := Timestampz(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, datetime.Location.String())
	endOfDay := Timestampz(
		now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, datetime.Location.String(),
	).SUB(INTERVALd(1 * time.Microsecond)) // Precision of Postgres is microsecond.

	stmt := SELECT(
		ClassGroupSessions.AllColumns,
	).FROM(
		ClassGroupSessions,
	).WHERE(
		ClassGroupSessions.StartTime.GT_EQ(startOfDay).AND(
			ClassGroupSessions.EndTime.LT_EQ(endOfDay),
		),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
