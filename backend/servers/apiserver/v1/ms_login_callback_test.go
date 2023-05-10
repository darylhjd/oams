package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_msLoginCallback(t *testing.T) {
	tests := []struct {
		name     string
		state    string
		httpCode int
	}{
		{
			"expected state, accepted callback",
			namespace,
			http.StatusOK,
		},
		{
			"unexpected state, rejecting callback",
			"wrong-state",
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

			postForm := url.Values{}
			postForm.Set(callbackStateParam, tt.state)
			postForm.Set(postFormCodeParam, "testing-code")

			resp, err := http.PostForm(reqUrl, postForm)
			a.Nil(err)

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			a.Equal(tt.httpCode, resp.StatusCode)
			if tt.httpCode != http.StatusOK {
				return
			}

			var actualBody msLoginCallbackResponse
			if err = json.Unmarshal(body, &actualBody); err != nil {
				t.Fatal(err)
			}

			a.Equal(msLoginCallbackResponse{AccessToken: "mock-access-token"}, actualBody)
		})
	}
}
