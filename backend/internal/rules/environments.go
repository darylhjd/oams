package rules

import "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"

type Environment interface {
	IsEnv() bool
}

type BaseEnvironment struct {
	Enrollments []model.SessionEnrollment `expr:"enrollments"`
}

func (e BaseEnvironment) IsEnv() bool {
	return true
}

type ConsecutiveEnvironment struct {
	BaseEnvironment
	ConsecutiveClasses int `expr:"consecutive_classes"`
}

type PercentageEnvironment struct {
	BaseEnvironment
	Percentage  float64 `expr:"percentage"`
	FromSession int     `expr:"from_session"`
}
