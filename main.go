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

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func ExecuteReq(query string, schema graphql.Schema) *graphql.Result {
	// fmt.Printf("Query: %v\n", query)
	// fmt.Printf("Schema: %v\n", schema)
	res := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	// fmt.Printf("Result: %v\n", res)
	if len(res.Errors) > 0 {
		fmt.Printf("Errors: %v\n", res.Errors)
	}
	return res
}

func main() {
	if SchemaError != nil {
		// Check for an error in schema at runtime
		panic(SchemaError)
	}

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	initRandom()
	ConnectDB()

	runDBTests()

	defer DisconnectDB()

	// router := mux.NewRouter()

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)

	// router.Handle("/graphql", h)

	srv := &http.Server{
		Handler: h,
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
	log.Println("\nShutting down")
	os.Exit(0)
}
