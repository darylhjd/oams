package webserver

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebServer_home(t *testing.T) {
	server := newTestWebServer()
	defer server.Close()

	reqUrl, err := url.JoinPath(server.URL, "/")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(reqUrl)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	a := assert.New(t)
	a.Equal(http.StatusOK, resp.StatusCode)
	a.Contains(string(body), "Welcome to Oats!")
}
