package v1

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (v *APIServerV1) msLoginCallback(w http.ResponseWriter, r *http.Request) {
	v.l.Debug(fmt.Sprintf("%s - received callback from azure", namespace), zap.String("method", r.Method))
}
