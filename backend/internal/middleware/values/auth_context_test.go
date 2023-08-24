package values

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthContext(t *testing.T) {
	tests := []struct {
		name         string
		contextValue any
		wantErr      error
	}{
		{
			"with proper context value",
			AuthContext{},
			nil,
		},
		{
			"bad context value",
			time.Time{},
			ErrUnexpectedAuthContextType,
		},
		{
			"no context value",
			nil,
			ErrNoAuthContext,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = req.WithContext(context.WithValue(req.Context(), AuthContextKey, tt.contextValue))

			switch {
			case tt.wantErr != nil:
				a.PanicsWithError(tt.wantErr.Error(), func() {
					GetAuthContext(req.Context())
				})
			default:
				a.NotPanics(func() {
					GetAuthContext(req.Context())
				})
			}
		})
	}
}
