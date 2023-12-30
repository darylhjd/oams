package environment

import (
	"github.com/darylhjd/oams/backend/internal/rules/environment/types"
)

type E interface {
	Type() types.T
}

type BaseE struct {
	EnvType     types.T  `json:"env_type"`
	Enrollments []string `expr:"enrollments" json:"-"`
}

func (e BaseE) Type() types.T {
	return e.EnvType
}

type ConsecutiveE struct {
	BaseE
	ConsecutiveClasses int `expr:"consecutive_classes" json:"consecutive_classes"`
}

type PercentageE struct {
	BaseE
	Percentage  float64 `expr:"percentage" json:"percentage"`
	FromSession int     `expr:"from_session" json:"from_session"`
}
