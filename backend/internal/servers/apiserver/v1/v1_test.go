package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/darylhjd/oams/backend/internal/middleware/values"
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/tests"
)

func newTestAPIServerV1(t *testing.T, dbId string) *APIServerV1 {
	t.Helper()

	testDb := tests.SetUp(t, dbId)
	return New(zap.NewNop(), testDb, tests.NewMockAzureAuthenticator())
}

func httpRequestWithAuthContext(r *http.Request, authContext values.AuthContext) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), values.AuthContextKey, authContext))
}
