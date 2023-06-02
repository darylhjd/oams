package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/env"
)

func TestAPIServerV1_login(t *testing.T) {
	tests := []struct {
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

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := newTestAPIServerV1(t)

			loginQueries := url.Values{}
			loginQueries.Set(stateReturnToQueryParam, tt.returnTo)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", loginUrl, loginQueries.Encode()), nil)
			rr := httptest.NewRecorder()
			v1.login(rr, req)

			var loginResp loginResponse
			err := json.Unmarshal(rr.Body.Bytes(), &loginResp)
			a.Nil(err)

			actualUrl, err := url.Parse(loginResp.RedirectUrl)
			a.Nil(err)

			s, err := json.Marshal(state{
				Version:  namespace,
				ReturnTo: tt.returnTo,
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

			a.Equal(http.StatusOK, rr.Code)
			a.Equal(url.URL{
				Scheme:   "https",
				Host:     "login.microsoftonline.com",
				Path:     path,
				RawQuery: expectedQueries.Encode(),
			}, *actualUrl)
		})
	}
}
