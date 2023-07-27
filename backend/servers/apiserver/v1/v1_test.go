package v1

import (
	"testing"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/tests"
)

func newTestAPIServerV1(t *testing.T, dbId string) *APIServerV1 {
	t.Helper()

	testDb := tests.SetUp(t, dbId)
	return New(zap.NewNop(), testDb, tests.NewMockAzureAuthenticator())
}

// ptr is a helper function to return a pointer to a type.
func ptr[T any](t T) *T {
	return &t
}
