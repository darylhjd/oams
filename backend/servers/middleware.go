package servers

import (
	"context"
	"log"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"

	"github.com/darylhjd/oams/backend/env"
)

const (
	ClaimsContextKey = "claims"
)

// AllowMethods allows a handler to accept only certain specified HTTP methods.
func AllowMethods(handlerFunc http.HandlerFunc, methods ...string) http.HandlerFunc {
	// Create set of allowed methods.
	allowedMethods := map[string]struct{}{}
	for _, method := range methods {
		allowedMethods[method] = struct{}{}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := allowedMethods[r.Method]; !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		handlerFunc(w, r)
	}
}

// CheckAuthorised checks if a request is authorised for a handler.
func CheckAuthorised(handlerFunc http.HandlerFunc, authenticator Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		set, err := authenticator.GetKeyCache().Get(r.Context(), keySetSource)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// We will be using session cookies for authentication.
		// TODO: Redirect user to login if required credentials are not present.
		c, err := r.Cookie(SessionCookieIdent)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		acct, err := authenticator.Account(r.Context(), c.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		res, err := authenticator.AcquireTokenSilent(
			r.Context(),
			[]string{env.GetAPIServerAzureLoginScope()},
			confidential.WithSilentAccount(acct))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, _, err := checkAzureToken(set, res.AccessToken)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Add the claims to the request context so other handlers/middleware can access it.
		r = r.WithContext(context.WithValue(r.Context(), ClaimsContextKey, claims))
		handlerFunc(w, r)
	}
}
