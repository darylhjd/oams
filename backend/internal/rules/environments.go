package rules

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/darylhjd/oams/backend/internal/intervention/fact"
)

// T is an identifier for the type of environment.
type T int

const (
	TConsecutive T = iota
	TPercentage
	TAdvanced
)

// Environment is a wrapper type for the database. A custom scanner and valuer is defined for this type.
type Environment struct {
	Env E `json:"env"`
}

func (e *Environment) Scan(value any) error {
	var t BaseE
	if err := json.Unmarshal(value.([]byte), &t); err != nil {
		return err
	}

	switch t.EnvType {
	case TConsecutive:
		e.Env = &ConsecutiveE{
			BaseE: BaseE{EnvType: t.EnvType},
		}
	case TPercentage:
		e.Env = &PercentageE{
			BaseE: BaseE{EnvType: t.EnvType},
		}
	case TAdvanced:
		e.Env = &BaseE{
			EnvType: t.EnvType,
		}
	default:
		return errors.New("unknown rule environment type")
	}

	return json.Unmarshal(value.([]byte), e.Env)
}

func (e *Environment) Value() (driver.Value, error) {
	return json.Marshal(e.Env)
}

// E is an interface that all environments for any rule must satisfy.
type E interface {
	SetFacts([]fact.F) E
}

// BaseE represents the base environment variables that all environments must have.
type BaseE struct {
	EnvType     T        `json:"env_type"`
	Enrollments []fact.F `expr:"enrollments" json:"-"`
}

func (e BaseE) SetFacts(facts []fact.F) E {
	e.Enrollments = facts
	return e
}

type ConsecutiveE struct {
	BaseE
	ConsecutiveClasses int `expr:"consecutive_classes" json:"consecutive_classes"`
}

func (e ConsecutiveE) SetFacts(facts []fact.F) E {
	e.Enrollments = facts
	return e
}

type PercentageE struct {
	BaseE
	Percentage  float64 `expr:"percentage" json:"percentage"`
	FromSession int     `expr:"from_session" json:"from_session"`
}

func (e PercentageE) SetFacts(facts []fact.F) E {
	e.Enrollments = facts
	return e
}