package v1

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/env"
)

const (
	postFormCodeParam = "code"
)

func (v *APIServerV1) msLoginCallback(w http.ResponseWriter, r *http.Request) {
	v.l.Debug(fmt.Sprintf("%s - received login callback from azure", namespace),
		zap.String("method", r.Method))

	// Check that we only handle callbacks from appropriate API version.
	state := r.PostFormValue(callbackStateParam)
	if state != namespace {
		w.WriteHeader(http.StatusTeapot)
		v.l.Error(fmt.Sprintf("%s - received login callback of different version so ignoring", namespace),
			zap.String(callbackStateParam, state))
		return
	}

	code := r.PostFormValue(postFormCodeParam)
	_, err := v.azure.AcquireTokenByAuthCode(
		r.Context(),
		code,
		env.GetAPIServerAzureLoginCallbackURL(),
		[]string{env.GetAPIServerAzureLoginScope()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not get tokens from auth code", namespace), zap.Error(err))
		return
	}

	response := "Authentication successful!"
	if _, err = w.Write([]byte(response)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("url", msLoginCallbackUrl),
			zap.Error(err))
	}
}
