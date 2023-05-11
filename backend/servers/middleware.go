package servers

import (
	"context"
	"log"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwk"
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
func CheckAuthorised(handlerFunc http.HandlerFunc, cache *jwk.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		set, err := cache.Get(r.Context(), keySetSource)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		claims, _, err := checkAzureToken(r, set)
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
