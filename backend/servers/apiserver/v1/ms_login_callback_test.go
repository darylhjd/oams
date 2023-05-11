package v1

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_msLoginCallback(t *testing.T) {
	tests := []struct {
		name     string
		state    state
		httpCode int
	}{
		{
			"expected state, accepted callback",
			state{
				Version: namespace,
			},
			http.StatusSeeOther,
		},
		{
			"expected state with custom return url, accepted callback",
			state{
				Version:  namespace,
				ReturnTo: "/randomUrl",
			},
			http.StatusSeeOther,
		},
		{
			"unexpected state, rejected callback",
			state{
				Version: "wrong-state",
			},
			http.StatusTeapot,
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

			expectedLocation := tt.state.ReturnTo
			if expectedLocation == "" {
				expectedLocation = Url
			}

			a.Equal(expectedLocation, resp.Header.Get("Location"))
		})
	}
}
