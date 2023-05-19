package tests

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/database"
)

func TestConnectDB(t *testing.T) {
	a := assert.New(t)
	a.Nil(testDb.C.Ping())
}

func TestGetConnectionProperties(t *testing.T) {
	a := assert.New(t)

	driver, connString := database.GetConnectionProperties(database.Namespace)

	u, err := url.Parse(connString)
	if err != nil {
		t.Fatal(err)
	}

	a.Equal("postgres", driver)
	a.Equal(driver, u.Scheme)
	a.Equal(fmt.Sprintf("/%s", database.Namespace), u.Path)
}