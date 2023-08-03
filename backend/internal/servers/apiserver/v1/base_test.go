package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_base(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name         string
		path         string
		wantResponse apiResponse
	}{
		{
			"valid request to base url",
			baseUrl,
			baseResponse{
				response: newSuccessResponse(),
				Message:  "Welcome to OAMS API Service V1! To get started, read the API docs.",
			},
		},
		{
			"invalid endpoint",
			"/bad-endpoint",
			newErrorResponse(http.StatusNotFound, "endpoint not implemented or is not supported"),
		},
	}

	for _, tt := range tts {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := assert.New(t)
			id := uuid.NewString()

			v1 := newTestAPIServerV1(t, id)
			defer tests.TearDown(t, v1.db, id)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()
			v1.base(rr, req)

			expectedBytes, err := json.Marshal(tt.wantResponse)
			a.Nil(err)
			a.Equal(string(expectedBytes), rr.Body.String())
		})
	}
}
