package servers

import (
	"context"
	"net/http"
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
func CheckAuthorised(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := checkAzureToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Provide the handler with the token in case it needs it.
		r = r.WithContext(context.WithValue(r.Context(), AuthFieldName, token))
		handlerFunc(w, r)
	}
}
