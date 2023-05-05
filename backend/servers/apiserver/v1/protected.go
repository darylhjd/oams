package v1

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (v *APIServerV1) protected(w http.ResponseWriter, r *http.Request) {
	v.l.Debug(fmt.Sprintf("%s - handling protected request", namespace))

	if _, err := w.Write([]byte("You are authenticated!")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("url", protectedUrl),
			zap.Error(err))
	}
}
