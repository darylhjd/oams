package v1

import (
	"net/http"

	"go.uber.org/zap"
)

// ping is the handler for the pingUrl.
func (v *APIServerV1) ping(w http.ResponseWriter, r *http.Request) {
	v.l.Debug("apiserver/v1 - handling ping request")

	response := `Pong~

Oats API Service is running normally!`

	if _, err := w.Write([]byte(response)); err != nil {
		v.l.Error("apiserver - could not write response",
			zap.String("url", pingUrl),
			zap.Error(err))
	}
}
