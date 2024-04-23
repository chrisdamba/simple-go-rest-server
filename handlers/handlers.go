package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/chrisdamba/simple-go-rest-server/db"
	"github.com/chrisdamba/simple-go-rest-server/models"
)

const (
	dbName = "getir-case-study"
	colName = "records"
)

// inMemoryDb is your in-memory store
var inMemoryDb = map[string]models.InMemoryRecord{}

// getNextID generates IDs for in-memory records
func getNextID() string {
	return fmt.Sprintf("id-%d", len(inMemoryDb))
}

// fetchFromMongo handles POST requests to fetch data from MongoDB.
func FetchFromMongo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Parse request payload
	var payload models.RequestPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI") // Retrieve connection string from environment
	if mongoURI == "" {
		fmt.Errorf("MONGO_URI environment variable not set")
	}
	client, err := db.ConnectToMongo(mongoURI)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	// Query MongoDB
	collection := client.Database(dbName).Collection(colName) 
	filter := bson.M{
		"createdAt": bson.M{
			"$gte": payload.StartDate,
			"$lte": payload.EndDate,
		},
		"totalCount": bson.M{
			"$gte": payload.MinCount,
			"$lte": payload.MaxCount,
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var records []models.MongoRecord
	if err = cursor.All(context.Background(), &records); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert records to the desired type
	var formattedRecords []struct {
		Key        string `json:"key,omitempty"`
		CreatedAt  string `json:"createdAt,omitempty"`
		TotalCount int    `json:"totalCount,omitempty"`
	}

	for _, record := range records {
		formattedRecords = append(formattedRecords, struct {
			Key        string `json:"key,omitempty"`
			CreatedAt  string `json:"createdAt,omitempty"`
			TotalCount int    `json:"totalCount,omitempty"`
		}{
			Key:        record.Key,
			CreatedAt:  record.CreatedAt.Format("2006-01-02T15:04:05Z"),
			TotalCount: record.TotalCount,
		})
	}

	// Formulate response
	responsePayload := models.ResponsePayload{
		Code:    0,
		Msg:     "Success",
		Records: formattedRecords,
	}
	jsonResponse, err := json.Marshal(responsePayload)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

