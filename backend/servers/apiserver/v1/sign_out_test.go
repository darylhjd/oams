package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/oauth2"
)

func Test_signOut(t *testing.T) {
	a := assert.New(t)
	v1 := newTestAPIServerV1(t)

	req := httptest.NewRequest(http.MethodGet, signOutUrl, nil)
	rr := httptest.NewRecorder()
	v1.signOut(rr, req)

	a.Equal(http.StatusOK, rr.Code)
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == oauth2.SessionCookieIdent {
			return
		}
	}
	a.FailNow("could not detect expected session cookie")
}
