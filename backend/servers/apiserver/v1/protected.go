package v1

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (v *APIServerV1) protected(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("You are authenticated!")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("url", protectedUrl),
			zap.Error(err))
	}
}
