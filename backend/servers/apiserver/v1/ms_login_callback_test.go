package v1

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/oauth2"
)

func TestAPIServerV1_msLoginCallback(t *testing.T) {
	tests := []struct {
		name        string
		state       state
		httpCode    int
		redirectUrl string
	}{
		{
			"expected state, accepted callback",
			state{
				Version: namespace,
			},
			http.StatusSeeOther,
			Url,
		},
		{
			"expected state with custom return url, accepted callback",
			state{
				Version:  namespace,
				ReturnTo: "/randomUrl",
			},
			http.StatusSeeOther,
			"/randomUrl",
		},
		{
			"unexpected state, rejected callback",
			state{
				Version: "wrong-state",
			},
			http.StatusTeapot,
			"",
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newTestAPIServerV1(t)
			defer server.Close()

			reqUrl, err := url.JoinPath(server.URL, msLoginCallbackUrl)
			if err != nil {
				t.Fatal(err)
			}

			s, err := json.Marshal(tt.state)
			if err != nil {
				t.Fatal(err)
			}

			postForm := url.Values{}
			postForm.Set(callbackStateParam, string(s))
			postForm.Set(postFormCodeParam, "testing-code")

			// Use custom client to prevent redirect.
			httpClient := http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := httpClient.PostForm(reqUrl, postForm)
			a.Nil(err)

			a.Equal(tt.httpCode, resp.StatusCode)
			if tt.httpCode != http.StatusSeeOther {
				return
			}

			a.Equal(tt.redirectUrl, resp.Header.Get("Location"))

			// Check that session cookie exists.
			for _, cookie := range resp.Cookies() {
				if cookie.Name == oauth2.SessionCookieIdent {
					return
				}
			}
			a.FailNow("could not detect expected session cookie")
		})
	}
}