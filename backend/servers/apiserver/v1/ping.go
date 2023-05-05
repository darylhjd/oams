package v1

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// ping is the handler for the pingUrl.
func (v *APIServerV1) ping(w http.ResponseWriter, r *http.Request) {
	v.l.Debug(fmt.Sprintf("%s - handling ping request", namespace))

	response := "Pong~\n\n" +
		"OAMS API Service is running normally!"

	if v.db.C.Ping() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response = "Uh oh, not connected!"
	}

	if _, err := w.Write([]byte(response)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("url", pingUrl),
			zap.Error(err))
	}
}
