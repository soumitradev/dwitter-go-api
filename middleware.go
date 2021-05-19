package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/graphql-go/handler"
)

func customMiddleware(next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ContextHandler(r.Context(), w, r)
	})
}

func LoggingHandler(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func ContentTypeHandler(next http.Handler) http.Handler {
	return handlers.ContentTypeHandler(next, "application/json", "application/graphql")
}

func RecoveryHandler(next http.Handler) http.Handler {
	return handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(next)
}
