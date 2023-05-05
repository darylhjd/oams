package v1

import (
	"net/http"
	"net/url"

	"github.com/darylhjd/oams/backend/env"
)

func (v *APIServerV1) login(w http.ResponseWriter, r *http.Request) {
	// https://learn.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow
	clientId, err := env.GetAPIServerAzureClientID()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	scope, err := env.GetAPIServerAzureLoginScope()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	redirectString, err := v.azure.AuthCodeURL(r.Context(), clientId,
		"http://localhost:8080/api/v1/ms-login-callback", []string{scope})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Add response_mode is form_post
	redirectUrl, err := url.Parse(redirectString)
	values := redirectUrl.Query()
	values.Set("response_mode", "form_post")
	redirectUrl.RawQuery = values.Encode()
	redirectString = redirectUrl.String()

	http.Redirect(w, r, redirectString, http.StatusSeeOther)
}
