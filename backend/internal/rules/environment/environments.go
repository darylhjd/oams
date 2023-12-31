package environment

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/darylhjd/oams/backend/internal/rules/environment/types"
	"github.com/darylhjd/oams/backend/internal/rules/fact"
)

// Environment is a wrapper type for the database. A custom scanner and valuer is defined for this type.
type Environment struct {
	Data E `json:"data"`
}

func (e *Environment) Scan(value any) error {
	var t BaseE
	if err := json.Unmarshal(value.([]byte), &t); err != nil {
		return err
	}

	switch t.EnvType {
	case types.TConsecutive:
		e.Data = &ConsecutiveE{
			BaseE: BaseE{EnvType: t.EnvType},
		}
	case types.TPercentage:
		e.Data = &PercentageE{
			BaseE: BaseE{EnvType: t.EnvType},
		}
	case types.TAdvanced:
		e.Data = &BaseE{
			EnvType: t.EnvType,
		}
	default:
		return errors.New("unknown rule environment type")
	}

	return json.Unmarshal(value.([]byte), e.Data)
}

func (e *Environment) Value() (driver.Value, error) {
	return json.Marshal(e.Data)
}

// E is an interface that all environments for any rule must satisfy.
type E interface {
	Type() types.T
}

// BaseE represents the base environment variables that all environments must have.
type BaseE struct {
	EnvType     types.T  `json:"env_type"`
	Enrollments []fact.F `expr:"enrollments" json:"-"`
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
