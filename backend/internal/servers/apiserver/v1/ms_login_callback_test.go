package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/darylhjd/oams/backend/pkg/to"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/oauth2"
)

func TestAPIServerV1_msLoginCallback(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name             string
		withState        state
		withExistingUser bool
		wantCode         int
		wantBody         apiResponse
		wantRedirectUrl  string
	}{
		{
			"expected state, accepted callback with new user",
			state{
				Version: namespace,
			},
			false,
			http.StatusSeeOther,
			nil,
			Url,
		},
		{
			"expected state, accepted callback with existing user",
			state{
				Version: namespace,
			},
			true,
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
			false,
			http.StatusSeeOther,
			nil,
			"/randomUrl",
		},
		{
			"unexpected version in state, rejected callback",
			state{
				Version: "wrong-version",
			},
			false,
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
			ctx := context.Background()
			id := uuid.NewString()

			v1 := newTestAPIServerV1(t, id)
			defer tests.TearDown(t, v1.db, id)

			var (
				expectedUser model.User
				err          error
			)
			if tt.withExistingUser {
				tests.StubAuthContextUser(t, ctx, v1.db)
				// Set a name to the auth context user, so we can check for no change.
				expectedUser, err = v1.db.UpdateUser(ctx, tests.MockAuthenticatorIDTokenName, database.UpdateUserParams{
					Name: to.Ptr("TEST ACCOUNT NAME LIM"),
				})
				a.Nil(err)
			}

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
			if tt.withExistingUser {
				// Check that if the user is an existing user, then user's information is unchanged in database.
				checkUser, err := v1.db.GetUser(ctx, tests.MockAuthenticatorIDTokenName)
				a.Nil(err)
				a.Equal(expectedUser, checkUser)
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
