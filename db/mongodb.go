package db

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectToMongo connects to MongoDB and returns a client.
func ConnectToMongo(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Failed to connect to MongoDB:", err)
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("Failed to ping MongoDB:", err)
		return nil, fmt.Errorf("error pinging MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")
	return client, nil
}


// fetchFromMongoDB handles POST requests to fetch data from MongoDB.
func fetchFromMongoDB(w http.ResponseWriter, r *http.Request) {
	// ...
}