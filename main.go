package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chrisdamba/simple-go-rest-server/handlers"
)

const (
	port        = 8080 // Default port
)

func main() {
	http.HandleFunc("/mongo", handlers.FetchFromMongo)
	http.HandleFunc("/in-memory", handlers.InMemoryHandler)

	// Start server
	log.Printf("Server listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
