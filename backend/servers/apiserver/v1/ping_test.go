package v1

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_ping(t *testing.T) {
	a := assert.New(t)

	server := newTestAPIServerV1(t)
	defer server.Close()

	reqUrl, err := url.JoinPath(server.URL, pingUrl)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(reqUrl)
	a.Nil(err)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	a.Equal(http.StatusOK, resp.StatusCode)
	a.Contains(string(body), "Pong~")
}
