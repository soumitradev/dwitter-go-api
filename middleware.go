package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func customMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > (65 << 20) {
			msg := "Request too large."
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LoggingHandler(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func ContentTypeHandler(next http.Handler) http.Handler {
	return handlers.ContentTypeHandler(next, "application/json", "application/graphql", "multipart/form-data")
}

func RecoveryHandler(next http.Handler) http.Handler {
	return handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(next)
}
