package v1

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

const (
	baseUrl = "/"
	pingUrl = "/ping"
)

type APIServerV1 struct {
	L *zap.Logger
}

func (v *APIServerV1) GetLogger() *zap.Logger {
	return v.L
}

// base is the handler for the baseUrl.
func (v *APIServerV1) base(w http.ResponseWriter, r *http.Request) {
	// Check if it is really the base URL.
	if r.URL.Path != baseUrl {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response :=
		`Welcome to Oats API Service V1!

To get started, read the API docs.`

	if _, err := w.Write([]byte(response)); err != nil {
		msg := fmt.Sprintf("apiserver - could not write response for %s: %s", baseUrl, err)
		v.L.Error(msg)
	}
}

// ping is the handler for the pingUrl.
func (v *APIServerV1) ping(w http.ResponseWriter, r *http.Request) {
	response := `Pong~

Oats API Service is running normally!`

	if _, err := w.Write([]byte(response)); err != nil {
		msg := fmt.Sprintf("apiserver - could not write response for %s: %s", pingUrl, err)
		v.L.Error(msg)
	}
}

func (v *APIServerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc(baseUrl, v.base)
	mux.HandleFunc(pingUrl, v.ping)

	mux.ServeHTTP(w, r)
}
