package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// Parse request payload
	var payload models.RequestPayload
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()  // This prevents decoding issues with unexpected fields
	if err := decoder.Decode(&payload); err != nil {
			respondWithError(w, 1, fmt.Sprintf("Invalid request payload: %v", err))
			return
	}
	// Perform further payload validation as needed
	if payload.StartDate == "" || payload.EndDate == "" {
		respondWithError(w, 2, "StartDate and EndDate fields are required")
		return
	}
	if _, err := time.Parse("2006-01-02", payload.StartDate); err != nil {
		respondWithError(w, 3, "Invalid StartDate format. Please use YYYY-MM-DD.")
		return
	}
	if _, err := time.Parse("2006-01-02", payload.EndDate); err != nil {
		respondWithError(w, 4, "Invalid EndDate format. Please use YYYY-MM-DD.")
		return
	}

	// Assume `minCount` and `maxCount` are required fields and check them
	if payload.MinCount < 0 || payload.MaxCount < 0 || payload.MinCount > payload.MaxCount {
		respondWithError(w, 5, "minCount and maxCount must be non-negative and minCount must not exceed maxCount")
		return
	}

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI") // Retrieve connection string from environment
	if mongoURI == "" {
    log.Println("MONGO_URI environment variable not set")
    respondWithError(w, 6, "MONGO_URI environment variable not set")
    return
	}
	client, err := db.ConnectToMongo(mongoURI)
	if err != nil {
    respondWithError(w, 7, "Error connecting to MongoDB")
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


// respondWithError sends an error response with a given code and message.
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest) // Set the status code to BadRequest for all validation errors
	errorResponse := models.ErrorResponse{
		Code:    code,
		Message: message,
	}
	json.NewEncoder(w).Encode(errorResponse)
}
