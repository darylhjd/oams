package v1

import (
	"net/http"

	"go.uber.org/zap"
)

// base is the handler for the baseUrl.
func (v *APIServerV1) base(w http.ResponseWriter, r *http.Request) {
	// Check if it is really the base URL.
	if r.URL.Path != baseUrl {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	v.l.Debug("apiserver/v1 - handling base request")

	response :=
		`Welcome to Oats API Service V1!

To get started, read the API docs.`

	if _, err := w.Write([]byte(response)); err != nil {
		v.l.Error("apiserver - could not write response",
			zap.String("url", baseUrl),
			zap.Error(err))
	}
}
