package database

import (
	"context"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

type CoordinatingClass struct {
	ID        int64  `alias:"class.id" json:"id"`
	Code      string `alias:"class.code" json:"code"`
	Year      int32  `alias:"class.year" json:"year"`
	Semester  string `alias:"class.semester" json:"semester"`
	Programme string `alias:"class.programme" json:"programme"`
	Au        int16  `alias:"class.au" json:"au"`
}

func (d *DB) GetCoordinatingClasses(ctx context.Context) ([]CoordinatingClass, error) {
	var res []CoordinatingClass

	stmt := selectCoordinatingClassFields().WHERE(
		coordinatingClassRLS(ctx),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) GetCoordinatingClass(ctx context.Context, id int64) (CoordinatingClass, error) {
	var res CoordinatingClass

	stmt := selectCoordinatingClassFields().WHERE(
		Classes.ID.EQ(Int64(id)).AND(
			coordinatingClassRLS(ctx),
		),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func (d *DB) GetCoordinatingClassRules(ctx context.Context, id int64) ([]model.ClassAttendanceRule, error) {
	var res []model.ClassAttendanceRule

	stmt := SELECT(
		ClassAttendanceRules.AllColumns,
	).FROM(
		ClassAttendanceRules,
	).WHERE(
		ClassAttendanceRules.ClassID.EQ(Int64(id)).AND(
			classAttendanceRuleRLS(ctx),
		),
	)

	err := stmt.QueryContext(ctx, d.qe, &res)
	return res, err
}

func selectCoordinatingClassFields() SelectStatement {
	return SELECT(
		Classes.ID,
		Classes.Code,
		Classes.Year,
		Classes.Semester,
		Classes.Programme,
		Classes.Au,
	).FROM(
		Classes,
	).ORDER_BY(
		Classes.Year.DESC(),
		Classes.Semester.DESC(),
	)
}
