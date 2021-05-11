package main

import (
	"net/http"

	"github.com/graphql-go/handler"
)

func customMiddleware(next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ContextHandler(r.Context(), w, r)
	})
}
