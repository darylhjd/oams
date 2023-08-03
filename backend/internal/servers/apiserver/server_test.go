package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/servers/apiserver/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAPIServer(t *testing.T) {
	tts := []struct {
		name           string
		withPath       string
		wantStatusCode int
	}{
		{
			"registered path",
			v1.Url,
			http.StatusOK,
		},
		{
			"unregistered path",
			"/unregistered/path",
			http.StatusNotFound,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			server := newTestAPIServer(t)
			server.registerHandlers()

			req := httptest.NewRequest(http.MethodGet, tt.withPath, nil)
			rr := httptest.NewRecorder()
			server.ServeHTTP(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
			if tt.wantStatusCode != http.StatusOK {
				a.Contains(rr.Body.String(), "malformed url path")
			}
		})
	}
}

func newTestAPIServer(t *testing.T) *APIServer {
	t.Helper()

	return &APIServer{
		l:   zap.NewNop(),
		mux: http.NewServeMux(),
		v1:  v1.New(zap.NewNop(), nil, nil),
	}
}
