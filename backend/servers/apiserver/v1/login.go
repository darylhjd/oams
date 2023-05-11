package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/env"
)

const (
	callbackMethodParam    = "response_mode"
	callbackMethodFormPost = "form_post"
	callbackStateParam     = "state"
)

type loginResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

func (v *APIServerV1) login(w http.ResponseWriter, r *http.Request) {
	v.l.Debug(fmt.Sprintf("%s - handling login request", namespace))

	// https://learn.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow
	redirectString, err := v.azure.AuthCodeURL(
		r.Context(),
		env.GetAPIServerAzureClientID(),
		env.GetAPIServerAzureLoginCallbackURL(),
		[]string{env.GetAPIServerAzureLoginScope()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - error creating azure redirect url", namespace), zap.Error(err))
		return
	}

	// Add extra parameters to the request.
	redirectUrl, err := url.Parse(redirectString)
	values := redirectUrl.Query()

	values.Set(callbackMethodParam, callbackMethodFormPost) // The callback is done through POST.
	values.Set(callbackStateParam, namespace)               // Include information on which API version the request originated from.

	redirectUrl.RawQuery = values.Encode()
	redirectString = redirectUrl.String()

	v.l.Debug(fmt.Sprintf("%s - generated azure login url", namespace), zap.String("url", redirectString))

	body, err := json.Marshal(loginResponse{RedirectUrl: redirectString})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - error marshalling url to body", namespace), zap.Error(err))
		return
	}

	if _, err = w.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("url", loginUrl),
			zap.Error(err))
	}
}
