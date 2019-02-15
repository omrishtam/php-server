package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"log"
)

func main() {
	address := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	userHandler := UserHandler{}
	r := mux.NewRouter()
	r.HandleFunc("/user/{id}", userHandler.GetUserHandler).Methods("GET")
	http.Handle("/", r)
	err := http.ListenAndServe(address + ":" + port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
