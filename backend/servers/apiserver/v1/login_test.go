package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/env"
)

func TestAPIServerV1_login(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name     string
		returnTo string
	}{
		{
			"login request with no return_to parameter",
			"",
		},
		{
			"login request with custom return_to parameter",
			"/return/to/here",
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

			loginQueries := url.Values{}
			loginQueries.Set(stateRedirectUrlQueryParam, tt.returnTo)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", loginUrl, loginQueries.Encode()), nil)
			rr := httptest.NewRecorder()
			v1.login(rr, req)

			s, err := json.Marshal(state{
				Version:     namespace,
				RedirectUrl: tt.returnTo,
			})
			a.Nil(err)

			expectedQueries := url.Values{}
			expectedQueries.Set("client_id", env.GetAPIServerAzureClientID())
			expectedQueries.Set("redirect_url", env.GetAPIServerAzureLoginCallbackURL())
			expectedQueries.Set("response_type", "code")
			expectedQueries.Set("scope", env.GetAPIServerAzureLoginScope())
			expectedQueries.Set(callbackMethodParam, callbackMethodFormPost)
			expectedQueries.Set(callbackStateParam, string(s))

			// NOTE: We add "/" to the beginning of the path so the test passes, but this will not affect the result.
			path, err := url.JoinPath("/", env.GetAPIServerAzureTenantID(), "oauth2", "v2.0", "authorize")
			a.Nil(err)

			expectedResp := loginResponse{
				newSuccessResponse(),
				(&url.URL{
					Scheme:   "https",
					Host:     "login.microsoftonline.com",
					Path:     path,
					RawQuery: expectedQueries.Encode(),
				}).String(),
			}

			expectedBytes, err := json.Marshal(expectedResp)
			a.Nil(err)
			a.Equal(string(expectedBytes), rr.Body.String())
		})
	}
}
