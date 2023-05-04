package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateArguments(t *testing.T) {
	tests := []struct {
		name        string
		args        arguments
		wantErr     bool
		containsErr string
	}{
		{
			"database name is empty",
			arguments{},
			true,
			"database name cannot be empty",
		},
		{
			"no operation specified",
			arguments{
				name: "testing",
			},
			true,
			"no operation specified",
		},
		{
			"migrate and create operations specified",
			arguments{
				name:    "testing",
				migrate: true,
				create:  true,
				drop:    false,
			},
			true,
			"more than one operation specified",
		},
		{
			"migrate and drop operations specified",
			arguments{
				name:    "testing",
				migrate: true,
				create:  false,
				drop:    true,
			},
			true,
			"more than one operation specified",
		},
		{
			"create and drop operations specified",
			arguments{
				name:    "testing",
				migrate: false,
				create:  true,
				drop:    true,
			},
			true,
			"more than one operation specified",
		},
		{
			"all operations specified",
			arguments{
				name:    "testing",
				migrate: true,
				create:  true,
				drop:    true,
			},
			true,
			"more than one operation specified",
		},
		{
			"migrate operation specified",
			arguments{
				name:    "testing",
				migrate: true,
				version: 1,
			},
			false,
			"",
		},
		{
			"no options specified for migrate operation",
			arguments{
				name:    "testing",
				migrate: true,
			},
			true,
			"no options specified",
		},
		{
			"more than one option specified for migrate operation",
			arguments{
				name:    "testing",
				migrate: true,
				version: 1,
				steps:   2,
			},
			true,
			"more than one option specified",
		},
		{
			"one option specified for migrate operation",
			arguments{
				name:    "testing",
				migrate: true,
				version: 1,
			},
			false,
			"",
		},
		{
			"create operation specified",
			arguments{
				name:   "testing",
				create: true,
			},
			false,
			"",
		},
		{
			"create operation specified but with options",
			arguments{
				name:   "testing",
				create: true,
				fullUp: true,
			},
			true,
			"invalid options provided for operation",
		},
		{
			"drop operation specified",
			arguments{
				name: "testing",
				drop: true,
			},
			false,
			"",
		},
		{
			"drop operation specified but with options",
			arguments{
				name:     "testing",
				drop:     true,
				fullDown: true,
			},
			true,
			"invalid options provided for operation",
		},
		{
			"truncate operation specified",
			arguments{
				name:     "testing",
				truncate: true,
			},
			false,
			"",
		},
		{
			"truncate operation specified but with options",
			arguments{
				name:     "testing",
				truncate: true,
				fullDown: true,
			},
			true,
			"invalid options provided for operation",
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateArguments(&(tt.args))
			a.Equal(tt.wantErr, err != nil)

			if tt.wantErr {
				a.ErrorContains(err, tt.containsErr)
			}
		})
	}
}
