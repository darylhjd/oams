package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const (
	stateRedirectUrlQueryParam = "redirect_url"
)

// oauthState stores the state from which the login was called.
// This helps us store useful information such as the redirect URL to return the user to after login
// and the version of the API in which the auth code flow was initiated.
type oauthState struct {
	Version     string `json:"version"`
	Verifier    string `json:"verifier"`
	RedirectUrl string `json:"redirect_url"`
}

type loginResponse struct {
	response
	RedirectUrl string `json:"redirect_url"`
}

func (v *APIServerV1) login(w http.ResponseWriter, r *http.Request) {
	verifier := oauth2.GenerateVerifier()
	// Set up the auth code flow state.
	s, err := json.Marshal(oauthState{
		Version:     namespace,
		Verifier:    verifier,
		RedirectUrl: r.URL.Query().Get(stateRedirectUrlQueryParam),
	})
	if err != nil {
		v.logInternalServerError(r, err)
		v.writeResponse(w, r, newErrorResponse(http.StatusInternalServerError, "cannot create oauth oauthState"))
		return
	}

	redirectString := v.auth.AuthCodeURL(string(s), verifier)

	v.l.Debug(fmt.Sprintf("%s - generated azure login url", namespace), zap.String("url", redirectString))
	v.writeResponse(w, r, loginResponse{
		newSuccessResponse(),
		redirectString,
	})
}
