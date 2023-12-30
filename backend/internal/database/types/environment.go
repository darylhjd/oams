package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/darylhjd/oams/backend/internal/rules"
)

type Environment struct {
	Env rules.Environment `json:"env"`
}

func (e *Environment) Scan(value any) error {
	return json.Unmarshal(value.([]byte), e.Env)
}

func (e *Environment) Value() (driver.Value, error) {
	return json.Marshal(e)
}
