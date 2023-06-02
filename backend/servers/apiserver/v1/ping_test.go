package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_ping(t *testing.T) {
	a := assert.New(t)
	v1 := newTestAPIServerV1(t)

	req := httptest.NewRequest(http.MethodGet, pingUrl, nil)
	rr := httptest.NewRecorder()
	v1.ping(rr, req)

	a.Equal(http.StatusOK, rr.Code)
	a.Contains(string(rr.Body.Bytes()), "Pong~")
}
