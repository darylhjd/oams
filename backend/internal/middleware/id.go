package middleware

import "net/http"

type HandlerWithID[T comparable] func(http.ResponseWriter, *http.Request, T)

// WithID allows usage of a handler that needs an ID identifier.
func WithID[T comparable](id T, handler HandlerWithID[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, id)
	}
}
