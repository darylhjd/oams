package v1

import (
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

	// Allows us to check if a redirect was done (and prevent it, so we can test the handler).
	httpClient := http.DefaultClient
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := httpClient.Get(reqUrl)
	a.Nil(err)

	redirect, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}

	expectedQueries := url.Values{}
	expectedQueries.Set("client_id", env.GetAPIServerAzureClientID())
	expectedQueries.Set("redirect_url", env.GetAPIServerAzureLoginCallbackURL())
	expectedQueries.Set("response_type", "code")
	expectedQueries.Set("scope", env.GetAPIServerAzureLoginScope())
	expectedQueries.Set(callbackMethodParam, callbackMethodFormPost)
	expectedQueries.Set(callbackStateParam, namespace)

	// NOTE: We add "/" to the beginning of the path so the test passes, but this will not affect the result.
	path, err := url.JoinPath("/", env.GetAPIServerAzureTenantID(), "oauth2", "v2.0", "authorize")
	if err != nil {
		t.Fatal(err)
	}

	a.Equal(http.StatusSeeOther, resp.StatusCode)
	a.Equal(url.URL{
		Scheme:   "https",
		Host:     "login.microsoftonline.com",
		Path:     path,
		RawQuery: expectedQueries.Encode(),
	}, *redirect)
}
