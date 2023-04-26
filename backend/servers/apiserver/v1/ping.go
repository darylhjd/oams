package v1

import (
	"fmt"
	"net/http"
)

// ping is the handler for the pingUrl.
func (v *APIServerV1) ping(w http.ResponseWriter, r *http.Request) {
	response := `Pong~

Oats API Service is running normally!`

	if _, err := w.Write([]byte(response)); err != nil {
		msg := fmt.Sprintf("apiserver - could not write response for %s: %s", pingUrl, err)
		v.l.Error(msg)
	}
}
