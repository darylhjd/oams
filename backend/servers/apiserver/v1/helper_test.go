package v1

import (
	"testing"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/tests"
)

func newTestAPIServerV1(t *testing.T, dbId string) *APIServerV1 {
	t.Helper()

	testDb := tests.SetUp(t, dbId)
	return NewAPIServerV1(zap.NewNop(), testDb, tests.NewMockAzureAuthenticator())
}
