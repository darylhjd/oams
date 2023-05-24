package webserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWebServer(t *testing.T) {
	server := newTestWebServer(t)
	defer server.Close()

	tests := []string{
		// For requests to pages (and not files), we expect StatusOK response.
		// Let the frontend router handle the file!
		"/",
		"/about",
		"/profile",
		"/nested/pages",
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(fmt.Sprintf("with request path %+q", tt), func(t *testing.T) {
			reqUrl, err := url.JoinPath(server.URL, tt)
			a.Nil(err)

			resp, err := http.Get(reqUrl)
			a.Nil(err)

			a.Equal(http.StatusOK, resp.StatusCode)
		})
	}
}

func newTestWebServer(t *testing.T) *httptest.Server {
	t.Helper()

	webServer, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Disable logger for tests.
	webServer.l = zap.NewNop()

	return httptest.NewServer(webServer)
}
