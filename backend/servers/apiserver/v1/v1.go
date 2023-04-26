package v1

import (
	"net/http"

	"go.uber.org/zap"
)

type APIServerV1 struct {
	L *zap.Logger
}

func (v *APIServerV1) GetLogger() *zap.Logger {
	return v.L
}

func (v *APIServerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.L.Info("v1 request being handled...")
	v.L.Info(r.URL.Path)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		v.L.Info("inside v1 base.")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("In v1 base!"))
	})

	mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		v.L.Info("inside hello")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Hello, World!"))
	})

	mux.ServeHTTP(w, r)
}
