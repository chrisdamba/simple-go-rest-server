package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


const dbName   = "getir-case-study"         
const colName  = "records" 

func connectToMongo() (*mongo.Client, error) {
	// Set client options
	mongoURI := os.Getenv("MONGODB_URI") // Retrieve connection string from environment
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable not set")
	}
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %v", err)
	}

	return client, nil
}


// fetchFromMongoDB handles POST requests to fetch data from MongoDB.
func fetchFromMongoDB(w http.ResponseWriter, r *http.Request) {
	// ...