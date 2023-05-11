package v1

import (
	"encoding/json"
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
	var s state
	err := json.Unmarshal([]byte(r.PostFormValue(callbackStateParam)), &s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - cannot parse state from login callback", namespace), zap.Error(err))
		return
	}

	if s.Version != namespace {
		w.WriteHeader(http.StatusTeapot)
		v.l.Error(fmt.Sprintf("%s - received login callback of different version so ignoring", namespace),
			zap.Any(callbackStateParam, s))
		return
	}

	// TODO: Find proper way to return tokens.
	_, err = v.azure.AcquireTokenByAuthCode(
		r.Context(),
		r.PostFormValue(postFormCodeParam),
		env.GetAPIServerAzureLoginCallbackURL(),
		[]string{env.GetAPIServerAzureLoginScope()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not get tokens from auth code", namespace), zap.Error(err))
		return
	}

	redirectUrl := s.ReturnTo
	if redirectUrl == "" {
		redirectUrl = Url
	}

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}
