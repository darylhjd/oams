package webserver

import (
	"fmt"
	"net/http"
)

const (
	homeUrl = "/"
)

func (s *WebServer) home(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Welcome to Oats!")); err != nil {
		msg := fmt.Sprintf("webserver - could not write response for %s: %s", homeUrl, err)
		s.l.Error(msg)
	}
}
