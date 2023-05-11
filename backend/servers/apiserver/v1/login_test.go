package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/env"
)

func TestAPIServerV1_login(t *testing.T) {
	a := assert.New(t)

	server := newTestAPIServerV1(t)
	defer server.Close()

	reqUrl, err := url.JoinPath(server.URL, loginUrl)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(reqUrl)
	a.Nil(err)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var loginResp loginResponse
	if err = json.Unmarshal(body, &loginResp); err != nil {
		t.Fatal(err)
	}

	actualUrl, err := url.Parse(loginResp.RedirectUrl)
	a.Nil(err)

	s, err := json.Marshal(state{
		Version: namespace,
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedQueries := url.Values{}
	expectedQueries.Set("client_id", env.GetAPIServerAzureClientID())
	expectedQueries.Set("redirect_url", env.GetAPIServerAzureLoginCallbackURL())
	expectedQueries.Set("response_type", "code")
	expectedQueries.Set("scope", env.GetAPIServerAzureLoginScope())
	expectedQueries.Set(callbackMethodParam, callbackMethodFormPost)
	expectedQueries.Set(callbackStateParam, string(s))

	// NOTE: We add "/" to the beginning of the path so the test passes, but this will not affect the result.
	path, err := url.JoinPath("/", env.GetAPIServerAzureTenantID(), "oauth2", "v2.0", "authorize")
	if err != nil {
		t.Fatal(err)
	}

	a.Equal(http.StatusOK, resp.StatusCode)
	a.Equal(url.URL{
		Scheme:   "https",
		Host:     "login.microsoftonline.com",
		Path:     path,
		RawQuery: expectedQueries.Encode(),
	}, *actualUrl)
}
