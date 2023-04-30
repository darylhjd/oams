package v1

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/database"
	"github.com/darylhjd/oams/backend/tests"
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

func newTestAPIServerV1(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(NewAPIServerV1(zap.NewNop(), testDb))
}
