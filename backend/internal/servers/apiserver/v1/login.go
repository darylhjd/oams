package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/env"
)

const (
	callbackMethodParam        = "response_mode"
	callbackMethodFormPost     = "form_post"
	callbackStateParam         = "state"
	stateRedirectUrlQueryParam = "redirect_url"
)

// state stores the state from which the login was called.
// This helps us store useful information such as the redirect URL to return the user to after login
// and the version of the API in which the auth code flow was initiated.
type state struct {
	Version     string `json:"version"`
	RedirectUrl string `json:"redirect_url"`
}

type loginResponse struct {
	response
	RedirectUrl string `json:"redirect_url"`
}

// login is the entrypoint to OAM's OAuth2 flow.
func (v *APIServerV1) login(w http.ResponseWriter, r *http.Request) {
	// https://learn.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow
	redirectString, err := v.azure.AuthCodeURL(
		r.Context(),
		env.GetAPIServerAzureClientID(),
		env.GetAPIServerAzureLoginCallbackURL(),
		[]string{env.GetAPIServerAzureLoginScope()})
	if err != nil {
		v.writeResponse(w, loginUrl, newErrorResponse(http.StatusInternalServerError, "cannot create auth code url"))
		return
	}

	// Set up the auth code flow state.
	s, err := json.Marshal(state{
		Version:     namespace,
		RedirectUrl: r.URL.Query().Get(stateRedirectUrlQueryParam),
	})
	if err != nil {
		v.writeResponse(w, loginUrl, newErrorResponse(http.StatusInternalServerError, "cannot create oauth state"))
		return
	}

	redirectUrl, err := url.Parse(redirectString)
	values := redirectUrl.Query()
	values.Set(callbackMethodParam, callbackMethodFormPost) // The callback is done through POST.
	values.Set(callbackStateParam, string(s))
	redirectUrl.RawQuery = values.Encode()
	redirectString = redirectUrl.String()

	v.l.Debug(fmt.Sprintf("%s - generated azure login url", namespace), zap.String("url", redirectString))

	v.writeResponse(w, loginUrl, loginResponse{
		newSuccessResponse(),
		redirectString,
	})
}
