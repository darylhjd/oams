package v1

import (
	"fmt"
	"net/http"
)

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
		v.l.Error(msg)
	}
}
