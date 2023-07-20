/*
Package tests is the testing package for the database. We use a different package for the tests since the
tests package requires helpers from the database package, which results in an import cycle.
*/
package tests

import (
	"os"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/tests"
)

var testDb *database.DB

func TestMain(m *testing.M) {
	testDb = tests.SetUp(m, database.Namespace)

	res := m.Run()
	tests.TearDown(m, testDb, database.Namespace)

	os.Exit(res)
}
