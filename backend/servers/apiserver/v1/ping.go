package v1

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// ping helps check the health of the API server. It also checks the database connection.
func (v *APIServerV1) ping(w http.ResponseWriter, r *http.Request) {
	response := "Pong~\n\n" +
		"OAMS API Service is running normally!"

	if err := v.db.C.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		response = "Uh oh, not connected!"
	}

	if _, err := w.Write([]byte(response)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("url", pingUrl),
			zap.Error(err))
	}
}
