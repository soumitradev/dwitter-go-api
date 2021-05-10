package main

import (
	"net/http"

	"github.com/graphql-go/handler"
)

type HTTPKey string

// HTTP is the struct used to inject the response writer and request http structs.
type HTTP struct {
	W *http.ResponseWriter
	R *http.Request
}

func httpHeaderMiddleware(next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ContextHandler(r.Context(), w, r)
	})
}
