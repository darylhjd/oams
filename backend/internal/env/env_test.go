package env

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkVarEmpty(t *testing.T) {
	tests := map[string]bool{
		"non-empty variable":                        false,
		"  non-empty variable with leading space":   false,
		"non-empty variable with trailing space   ": false,
		"   ": true,
		"":    true,
	}

	a := assert.New(t)
	for tt, wantErr := range tests {
		t.Run(fmt.Sprintf("with variable %+q", tt), func(t *testing.T) {
			a.Equal(wantErr, checkVarEmpty(tt) != nil)
		})
	}
}
