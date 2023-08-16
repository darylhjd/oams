package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/oauth2"
)

func TestAPIServerV1_msLoginCallback(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name            string
		withState       state
		wantCode        int
		wantBody        apiResponse
		wantRedirectUrl string
	}{
		{
			"expected state, accepted callback",
			state{
				Version: namespace,
			},
			http.StatusSeeOther,
			nil,
			Url,
		},
		{
			"expected state with custom return url, accepted callback",
			state{
				Version:     namespace,
				RedirectUrl: "/randomUrl",
			},
			http.StatusSeeOther,
			nil,
			"/randomUrl",
		},
		{
			"unexpected version in state, rejected callback",
			state{
				Version: "wrong-version",
			},
			http.StatusTeapot,
			newErrorResponse(http.StatusInternalServerError, "wrong api version handling"),
			"",
		},
	}

	for _, tt := range tts {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := assert.New(t)
			id := uuid.NewString()

			v1 := newTestAPIServerV1(t, id)
			defer tests.TearDown(t, v1.db, id)

			stateBytes, err := json.Marshal(tt.withState)
			a.Nil(err)

			postForm := url.Values{}
			postForm.Set(callbackStateParam, string(stateBytes))
			postForm.Set(postFormCodeParam, uuid.NewString())
			req := httptest.NewRequest(http.MethodPost, msLoginCallbackUrl, strings.NewReader(postForm.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			v1.msLoginCallback(rr, req)

			a.Equal(tt.wantCode, rr.Code)
			a.Equal(tt.wantRedirectUrl, rr.Header().Get("Location"))

			if tt.wantCode != http.StatusSeeOther {
				expectedBytes, err := json.Marshal(tt.wantBody)
				a.Nil(err)
				a.Equal(string(expectedBytes), rr.Body.String())
				return
			}

			// Check that user that logged in exists in the database.
			tests.CheckUserExists(a, context.Background(), v1.db, tests.MockAuthenticatorIDTokenName)

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
