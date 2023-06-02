package v1

import (
	"log"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	"github.com/darylhjd/oams/backend/internal/tests"
)

// testDb provides a package-level sand-boxed database for testing.
var testDb *database.DB

func TestMain(m *testing.M) {
	testEnv, err := tests.SetUp(m, namespace)
	if err != nil {
		log.Fatal(err)
	}

	testDb = testEnv.Db
	res := m.Run()
	tests.TearDown(m, testEnv, namespace)

	os.Exit(res)
}

func newTestAPIServerV1(t *testing.T) *APIServerV1 {
	t.Helper()
	return NewAPIServerV1(zap.NewNop(), testDb, oauth2.NewMockAzureAuthenticator())
}
