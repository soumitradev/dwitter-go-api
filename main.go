package main

import (
	"context"
	"dwitter_go_graphql/auth"
	"dwitter_go_graphql/cdn"
	"dwitter_go_graphql/consts"
	"dwitter_go_graphql/database"
	"dwitter_go_graphql/gql"
	"dwitter_go_graphql/middleware"
	"dwitter_go_graphql/util"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/functionalfoundry/graphqlws"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	"github.com/unrolled/secure"
)

func main() {

	if gql.SchemaError != nil {
		// Check for an error in schema at runtime
		panic(gql.SchemaError)
	}

	// Set flag for timeout to close all connections before quitting
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Load .env
	godotenv.Load()

	// Seed the random function
	util.InitRandom()

	// Connect to database, and seed the database
	database.ConnectDB()
	cdn.InitCDN()
	database.RunDBTests()

	// When returning from main(), make sure to disconnect from database
	defer database.DisconnectDB()

	// Create a new router, and add middleware
	router := mux.NewRouter().StrictSlash(true)

	// Create a graphql query handler
	h := handler.New(&handler.Config{
		Schema:     &gql.Schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
		// This is a way to pass context about the request into the resolver function of graphql
		RootObjectFn: func(myCtx context.Context, r *http.Request) map[string]interface{} {
			// Pass down the authorization token to the graphql query
			authHeader := r.Header.Get("authorization")
			tokenString := auth.SplitAuthToken(authHeader)
			return map[string]interface{}{
				"token": tokenString,
			}
		},
	})

	consts.SubscriptionManager = graphqlws.NewSubscriptionManager(&gql.Schema)

	graphqlwsHandler := graphqlws.NewHandler(graphqlws.HandlerConfig{
		// Wire up the GraphqL WebSocket handler with the subscription manager
		SubscriptionManager: consts.SubscriptionManager,

		// Optional: Add a hook to resolve auth tokens into users that are
		// then stored on the GraphQL WS connections
		Authenticate: func(authToken string) (interface{}, error) {
			data, _, err := auth.VerifyAccessToken(authToken)
			if err != nil {
				return nil, err
			}
			return data["username"].(string), nil
		},
	})

	// Map /graphql to the graphql handler, and attach a middleware to it
	router.Handle("/graphql", h)

	// Handle login using a non-GraphQL solution
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	router.HandleFunc("/refresh_token", auth.RefreshHandler).Methods("POST")
	router.HandleFunc("/media_upload", cdn.UploadMedia).Methods("POST")
	router.HandleFunc("/pfp_upload", cdn.UploadPfp).Methods("POST")
	router.Handle("/subscriptions", graphqlwsHandler)

	secureMiddleware := secure.New(secure.Options{
		FrameDeny: true,
	})

	router.Use(handlers.CompressHandler)
	router.Use(middleware.LoggingHandler)
	router.Use(middleware.ContentTypeHandler)
	router.Use(middleware.RecoveryHandler)
	router.Use(middleware.CustomMiddleware)
	router.Use(secureMiddleware.Handler)

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

	gql.InitSubscriptions()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	BaseCtx, cancel := context.WithTimeout(consts.BaseCtx, wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(BaseCtx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-main.BaseCtx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("Shutting down")
	os.Exit(0)
}
