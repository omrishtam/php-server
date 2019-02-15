package main

import (
	"log"
	"net/http"
	"os"
	"strings"
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

	err := http.ListenAndServe(address+":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("Home Page"))
	}
}

func configDB(ctx context.Context) (*mongo.Database, error) {
	uri := fmt.Sprintf(`mongodb://%s/%s`,
		ctx.Value(usernameKey),
		ctx.Value(databaseKey),
	)
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
