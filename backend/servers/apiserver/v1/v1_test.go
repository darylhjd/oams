package v1

import (
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func newTestAPIServerV1NoDB(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(NewAPIServerV1(zap.NewNop(), nil))
}
