package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/oauth2"
)

func TestAPIServerV1_logOut(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	ctx := context.Background()
	id := uuid.NewString()

	v1 := newTestAPIServerV1(t, id)
	defer tests.TearDown(t, v1.db, id)

	tests.StubAuthContextUser(t, ctx, v1.db)

	req := httpRequestWithAuthContext(
		httptest.NewRequest(http.MethodGet, logoutUrl, nil),
		tests.StubAuthContext(),
	)
	rr := httptest.NewRecorder()
	v1.logout(rr, req)

	expectedBytes, err := json.Marshal(logoutResponse{newSuccessResponse()})
	a.Nil(err)
	a.Equal(string(expectedBytes), rr.Body.String())

	// Check for session deletion cookie.
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == oauth2.SessionCookieIdent {
			return
		}
	}
	a.FailNow("could not detect expected session deletion cookie")
}
