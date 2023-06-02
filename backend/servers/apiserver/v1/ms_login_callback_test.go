package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/oauth2"
)

func TestAPIServerV1_msLoginCallback(t *testing.T) {
	tests := []struct {
		name            string
		state           state
		wantCode        int
		wantBody        string
		wantRedirectUrl string
	}{
		{
			"expected state, accepted callback",
			state{
				Version: namespace,
			},
			http.StatusSeeOther,
			"",
			Url,
		},
		{
			"expected state with custom return url, accepted callback",
			state{
				Version:  namespace,
				ReturnTo: "/randomUrl",
			},
			http.StatusSeeOther,
			"",
			"/randomUrl",
		},
		{
			"unexpected version in state, rejected callback",
			state{
				Version: "wrong-state",
			},
			http.StatusTeapot,
			"wrong API version used for login callback",
			"",
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := newTestAPIServerV1(t)

			s, err := json.Marshal(tt.state)
			a.Nil(err)

			postForm := url.Values{}
			postForm.Set(callbackStateParam, string(s))
			postForm.Set(postFormCodeParam, "testing-code")
			req := httptest.NewRequest(http.MethodPost, msLoginCallbackUrl, strings.NewReader(postForm.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			v1.msLoginCallback(rr, req)

			a.Equal(tt.wantCode, rr.Code)
			a.Equal(tt.wantRedirectUrl, rr.Header().Get("Location"))
			a.Contains(string(rr.Body.Bytes()), tt.wantBody)
			if tt.wantCode != http.StatusSeeOther {
				return
			}

			// Check that session cookie exists.
			for _, cookie := range rr.Result().Cookies() {
				if cookie.Name == oauth2.SessionCookieIdent {
					return
				}
			}
			a.FailNow("could not detect expected session cookie")
		})
	}
}
