package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
)

// func ExecuteReq(query string, schema graphql.Schema) *graphql.Result {
// 	// ctx := context.WithValue(context.Background(), "token", request.URL.Query().Get("token"))
// 	res := graphql.Do(graphql.Params{
// 		Schema:        schema,
// 		RequestString: query,
// 		Context:       ctx,
// 	})

// 	if len(res.Errors) > 0 {
// 		fmt.Printf("Errors: %v\n", res.Errors)
// 	}
// 	return res
// }

func main() {
	if SchemaError != nil {
		// Check for an error in schema at runtime
		panic(SchemaError)
	}

	// Set flag for timeout to close all connections before quitting
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Seed the random function
	initRandom()

	// Connect to database, and seed the database
	ConnectDB()
	runDBTests()

	// When returning from main(), make sure to disconnect from database
	defer DisconnectDB()

	// Create a new router
	router := mux.NewRouter()

	// Create a graphql query handler
	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
		// This is a way to pass context about the request into the resolver function of graphql
		RootObjectFn: func(myCtx context.Context, r *http.Request) map[string]interface{} {
			// Pass down the authorization token to the graphql query
			auth := r.Header.Get("authorization")
			tokenString := SplitAuthToken(auth)
			return map[string]interface{}{
				"token": tokenString,
			}
		},
	})

	// Map /graphql to the graphql handler, and attach a middleware to it
	router.Handle("/graphql", customMiddleware(h))

	// Handle login using a non-GraphQL solution
	router.HandleFunc("/login", loginHandler).Methods("POST")

	// Create an HTTP server
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:5000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	fmt.Println("Server now running on port 5000, access /graphql")

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println()
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("Shutting down")
	os.Exit(0)
}
