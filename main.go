package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"log"
)

const (
	hostKey     = key("hostKey")
	databaseKey = key("databaseKey")
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
	ctx = context.WithValue(ctx, databaseKey, os.Getenv("MONGO_DATABASE"))
	db, err := configDB(ctx)
	if err != nil {
		log.Fatalf("database configuration failed: %v", err)
	}
	
	userHandler := UserHandler{}
	r := mux.NewRouter()
	r.HandleFunc("/user/{id}", userHandler.GetUserHandler).Methods("GET")
	http.Handle("/", r)	
	err := http.ListenAndServe(address+":"+port, nil)
	if err != nil {
