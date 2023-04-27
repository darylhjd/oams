package webserver

import "net/http/httptest"

func newTestWebServer() *httptest.Server {
	webServer, _ := New()
	return httptest.NewServer(webServer)
}
