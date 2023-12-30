package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/darylhjd/oams/backend/internal/rules/environment"
	"github.com/darylhjd/oams/backend/internal/rules/environment/types"
)

type Environment struct {
	Env environment.E `json:"env"`
}

func (e *Environment) Scan(value any) error {
	t := struct {
		EnvType types.T `json:"env_type"`
	}{}
	if err := json.Unmarshal(value.([]byte), &t); err != nil {
		return err
	}

	switch t.EnvType {
	case types.TConsecutive:
		e.Env = &environment.ConsecutiveE{
			BaseE: environment.BaseE{EnvType: t.EnvType},
		}
	case types.TPercentage:
		e.Env = &environment.PercentageE{
			BaseE: environment.BaseE{EnvType: t.EnvType},
		}
	case types.TAdvanced:
		e.Env = &environment.BaseE{
			EnvType: t.EnvType,
		}
	}

	return json.Unmarshal(value.([]byte), e.Env)
}

func (e *Environment) Value() (driver.Value, error) {
	return json.Marshal(e)
}
