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

func TestAPIServerV1_ping(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	id := uuid.NewString()

	v1 := newTestAPIServerV1(t, id)
	defer tests.TearDown(t, v1.db, id)

	req := httptest.NewRequest(http.MethodGet, pingUrl, nil)
	rr := httptest.NewRecorder()
	v1.ping(rr, req)

	expectedBytes, err := json.Marshal(pingResponse{
		response: newSuccessResponse(),
		Message:  "Pong~ OAMS API Service is running normally!",
	})
	a.Nil(err)
	a.Equal(string(expectedBytes), rr.Body.String())
}
