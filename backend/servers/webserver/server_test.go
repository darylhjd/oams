package webserver

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWebServer(t *testing.T) {
	a := assert.New(t)

	server := newTestWebServer(t)
	defer server.Close()

	reqUrl, err := url.JoinPath(server.URL, "/")
	a.Nil(err)

	resp, err := http.Get(reqUrl)
	a.Nil(err)

	a.Equal(http.StatusOK, resp.StatusCode)
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
