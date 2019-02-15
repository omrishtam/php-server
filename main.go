package main

import (
	"fmt"
	"context"
	"log"
	"time"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/gorilla/mux"
)

type key string

const (
	hostKey     = key("hostKey")
	databaseKey = key("databaseKey")
	usernameKey = key("usernameKey")
	passwordKey = key("passwordKey")
)

func main() {
	var wait time.Duration
    flag.DurationVar(&wait,
		"graceful-timeout",
		time.Second * 15,
		"the duration for which the server gracefully" +
		"wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	
	address := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = context.WithValue(ctx, hostKey, os.Getenv("MONGO_HOST"))
	ctx = context.WithValue(ctx, usernameKey, os.Getenv("MONGO_USERNAME"))
	ctx = context.WithValue(ctx, passwordKey, os.Getenv("MONGO_PASSWORD"))
	ctx = context.WithValue(ctx, databaseKey, os.Getenv("MONGO_DATABASE"))
	db, err := configDB(ctx)
	if err != nil {
		log.Fatalf("database configuration failed: %v", err)
	}
	
	userHandler := UserHandler{}
	r := mux.NewRouter()
	r.HandleFunc("/user/{id}", userHandler.GetUserHandler).Methods("GET")
	r.HandleFunc("/user", userHandler.AddUserHandler).Methods("POST")
	r.HandleFunc("/user/{id}", userHandler.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", userHandler.DeleteUserHandler).Methods("DELETE")

	server := &http.Server{
        Addr:         address + ":" + port,
        // Good practice to set timeouts to avoid Slowloris attacks.
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: r, // Pass our instance of gorilla/mux in.
    }

	// Run our server in a goroutine so that it doesn't block.
    go func() {
        if err := server.ListenAndServe(); err != nil {
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
    gracefulShutdownCtx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    // Doesn't block if no connections, but will otherwise wait
    // until the timeout deadline.
    server.Shutdown(gracefulShutdownCtx)
    // Optionally, you could run srv.Shutdown in a goroutine and block on
    // <-ctx.Done() if your application should wait for other services
    // to finalize based on context cancellation.
    log.Println("shutting down")
    os.Exit(0)
}

func configDB(ctx context.Context) (*mongo.Database, error) {
	var uri string
	if ctx.Value(usernameKey) != nil && ctx.Value(passwordKey) != nil {
		uri = fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
			ctx.Value(usernameKey),
			ctx.Value(passwordKey),
			ctx.Value(hostKey),
			ctx.Value(databaseKey),
		)
	} else {
		uri = fmt.Sprintf(`mongodb://%s/%s`,
			ctx.Value(hostKey),
			ctx.Value(databaseKey),
		)
	}
	client, err := mongo.NewClient(uri)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to mongo: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("mongo client couldn't connect with background context: %v", err)
	}
	phpDB := client.Database("php")
	return phpDB, nil
}
