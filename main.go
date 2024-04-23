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

// main function
func main() {
	// Define endpoints and handlers
	http.HandleFunc("/mongo", handlers.FetchFromMongo)
	// http.HandleFunc("/in-memory", inMemoryHandler)

	// Start server
	log.Printf("Server listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
