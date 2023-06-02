package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_base(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantCode int
		wantBody string
	}{
		{
			"valid request to base url",
			baseUrl,
			http.StatusOK,
			"Welcome to OAMS API Service V1!",
		},
		{
			"invalid endpoint",
			"/bad-endpoint",
			http.StatusNotFound,
			"Endpoint not implemented or is not supported!",
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := newTestAPIServerV1(t)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()
			v1.base(rr, req)

			a.Equal(tt.wantCode, rr.Code)
			a.Contains(string(rr.Body.Bytes()), tt.wantBody)
		})
	}
}
