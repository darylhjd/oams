package tests

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/database"
)

func TestConnectDB(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	id := uuid.NewString()

	db := tests.SetUp(t, id)
	defer tests.TearDown(t, db, id)

	a.Nil(db.C.Ping(context.Background()))
}

func TestGetConnectionProperties(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	id := uuid.NewString()

	db := tests.SetUp(t, id)
	defer tests.TearDown(t, db, id)

	driver, connString := database.GetConnectionProperties(id)

	u, err := url.Parse(connString)
	if err != nil {
		t.Fatal(err)
	}

	a.Equal("postgres", driver)
	a.Equal(driver, u.Scheme)
	a.Equal(fmt.Sprintf("/%s", id), u.Path)
}
