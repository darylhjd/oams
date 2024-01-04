package rules

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Environment is a wrapper type for a database type. A custom scanner and valuer is defined for this type.
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

func (e Environment) Value() (driver.Value, error) {
	return json.Marshal(e.Env)
}

func (e *Environment) UnmarshalJSON(data []byte) error {
	return e.Scan(data)
}

func (e Environment) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Env)
}
