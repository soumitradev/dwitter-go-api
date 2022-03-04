// Package middleware provides useful custom middleware
package middleware

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/soumitradev/Dwitter/backend/common"

	"github.com/gorilla/handlers"
)

// Limit size of request
func SizeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > (65 << 20) {
			msg := "Request too large."
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Log requests
func LoggingHandler(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

// Limit content types
func ContentTypeHandler(next http.Handler) http.Handler {
	return handlers.ContentTypeHandler(next, "application/json", "application/graphql", "multipart/form-data")
}

// Handle recoveries
func RecoveryHandler(next http.Handler) http.Handler {
	return handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(next)
}

func CORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}
