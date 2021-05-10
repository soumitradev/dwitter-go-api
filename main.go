package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
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

type ContextKey string

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func main() {

	// re, err := regexp.Compile(`\{"q`)
	re, err := regexp.Compile(`(mutation\{login\(username:"|",password:"|"\)\{AccessToken\}\})`)
	if err != nil {
		fmt.Print(err)
		return
	}

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

	router := mux.NewRouter()

	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
		RootObjectFn: func(myCtx context.Context, r *http.Request) map[string]interface{} {
			myCtx = context.WithValue(myCtx, ContextKey("token"), r.URL.Query().Get("token"))
			return map[string]interface{}{
				"context": myCtx,
			}
		},
	})

	// router.Handle("/graphql", httpHeaderMiddleware(h))

	router.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		strBody := string(body)

		var result map[string]string
		json.Unmarshal([]byte(strBody), &result)

		queryText := result["query"]

		regexRes := re.Split(queryText, -1)

		user := regexRes[1]
		pass := regexRes[2]

		auth, _ := CheckCreds(user, pass)
		if auth {
			refTok, _ := RefreshToken(user)
			c := http.Cookie{
				Name:     "jid",
				Value:    refTok,
				HttpOnly: true,
				Secure:   true,
			}
			http.SetCookie(w, &c)
		}

		// And now set a new body, which will simulate the same data we read:
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// Create a response wrapper:
		mrw := &MyResponseWriter{
			ResponseWriter: w,
			buf:            &bytes.Buffer{},
		}

		h.ContextHandler(req.Context(), mrw, req)
		// req.Context()
	})

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
