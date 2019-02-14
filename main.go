package main

import (
	"net/http"
	"log"
)

const (
	address string = ":8080"
)

func main() {
	http.HandleFunc("/", homeHandler)
	log.Printf("Server listening on %s", address)
	http.ListenAndServe(address, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("Home Page"))
	}
} 
