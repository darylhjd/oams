package webserver

import "net/http/httptest"

func newTestWebServer() *httptest.Server {
	webServer, _ := NewWebServer()
	return httptest.NewServer(webServer)
}
