package main

import (
	"net/http"
	"os"
	"strings"
	"log"
)

func main() {
	envMap := mapEnv(os.Environ())

	http.HandleFunc("/", homeHandler)
	address := envMap["ADDRESS"]
	port := envMap["PORT"]
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(address + ":" + port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func mapEnv(env []string) map[string]string {
	mapped := map[string]string{}

	for _, val := range env {
		sep := strings.Index(val, "=")

		if sep > 0 {
			key := val[0:sep]
			value := val[sep+1:]
			mapped[key] = value
		} else {
			log.Println("Bad environment: " + val)
		}
	}

	return mapped
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("Home Page"))
	}
} 
