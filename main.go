package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"log"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"fmt"
)

type key string

const (
	hostKey     = key("hostKey")
	databaseKey = key("databaseKey")
	usernameKey = key("usernameKey")
	passwordKey = key("passwordKey")
)

func main() {
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
	http.Handle("/", r)	
	err = http.ListenAndServe(address+":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
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