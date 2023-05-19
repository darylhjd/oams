package v1

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// base is the handler for the baseUrl.
func (v *APIServerV1) base(w http.ResponseWriter, r *http.Request) {
	// Check if it is really the base URL.
	if r.URL.Path != baseUrl {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	v.l.Debug(fmt.Sprintf("%s - handling base request", namespace))

	response := "Welcome to OAMS API Service V1!\n\n" +
		"To get started, read the API docs."

	if _, err := w.Write([]byte(response)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("url", baseUrl),
			zap.Error(err))
	}
}
