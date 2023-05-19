package tests

import (
	"log"
	"os"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/tests"
)

var testDb *database.DB

func TestMain(m *testing.M) {
	testEnv, err := tests.SetUp(m, database.Namespace)
	if err != nil {
		log.Fatal(err)
	}

	testDb = testEnv.Db
	res := m.Run()
	tests.TearDown(m, testEnv, database.Namespace)

	os.Exit(res)
}
