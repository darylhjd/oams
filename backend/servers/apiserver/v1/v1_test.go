package v1

import (
	"net/http/httptest"

	"go.uber.org/zap"
)

func newTestAPIServerV1() *httptest.Server {
	return httptest.NewServer(NewAPIServerV1(zap.NewNop()))
}
