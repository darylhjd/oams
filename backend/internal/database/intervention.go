package database

import (
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/internal/intervention/fact"
	. "github.com/go-jet/jet/v2/postgres"
)

func GetIntervention() error {
	var _ fact.F
	var _ model.SessionEnrollment

	_ = SELECT(
		ClassAttendanceRules.AllColumns,
	).FROM(
		ClassAttendanceRules,
	)

	return nil
}
