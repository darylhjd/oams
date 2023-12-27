package database

import (
	"context"

	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

type CoordinatingClasses struct {
	ID        int64  `alias:"class.id" json:"id"`
	Code      string `alias:"class.code" json:"code"`
	Year      int32  `alias:"class.year" json:"year"`
	Semester  string `alias:"class.semester" json:"semester"`
	Programme string `alias:"class.programme" json:"programme"`
	Au        int16  `alias:"class.au" json:"au"`
}

func (d *DB) GetCoordinatingClasses(ctx context.Context) ([]CoordinatingClasses, error) {
	var res []CoordinatingClasses

	stmt := SELECT(
		Classes.ID,
		Classes.Code,
		Classes.Year,
		Classes.Semester,
		Classes.Programme,
		Classes.Au,
	).DISTINCT().FROM(
		Classes.INNER_JOIN(
			ClassGroups, ClassGroups.ClassID.EQ(Classes.ID),
		).LEFT_JOIN(
			ClassGroupManagers, ClassGroupManagers.ClassGroupID.EQ(ClassGroups.ID),
		),
	).WHERE(
		attendanceRuleRLS(ctx),
	).ORDER_BY(
		Classes.Year.DESC(),
		Classes.Semester.DESC(),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}
