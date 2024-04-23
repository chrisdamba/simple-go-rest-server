package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// InMemoryDb is the in-memory store
var InMemoryDb = map[string]models.InMemoryRecord{}

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
			respondWithError(w, 1, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
			return
	}
	// Perform further payload validation as needed
	if payload.StartDate == "" || payload.EndDate == "" {
		respondWithError(w, 2, http.StatusBadRequest, "StartDate and EndDate fields are required")
		return
	}
	if _, err := time.Parse("2006-01-02", payload.StartDate); err != nil {
		respondWithError(w, 3, http.StatusBadRequest, "Invalid StartDate format. Please use YYYY-MM-DD.")
		return
	}
	if _, err := time.Parse("2006-01-02", payload.EndDate); err != nil {
		respondWithError(w, 4, http.StatusBadRequest, "Invalid EndDate format. Please use YYYY-MM-DD.")
		return
	}

	// Assume `minCount` and `maxCount` are required fields and check them
	if payload.MinCount < 0 || payload.MaxCount < 0 || payload.MinCount > payload.MaxCount {
		respondWithError(w, 5, http.StatusBadRequest, "minCount and maxCount must be non-negative and minCount must not exceed maxCount")
		return
	}

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI") // Retrieve connection string from environment
	if mongoURI == "" {
    log.Println("MONGO_URI environment variable not set")
    respondWithError(w, 6, http.StatusInternalServerError, "MONGO_URI environment variable not set")
    return
	}
	client, err := db.ConnectToMongo(mongoURI)
	if err != nil {
    respondWithError(w, 7, http.StatusInternalServerError, "Error connecting to MongoDB")
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
	respondWithJSON(w, http.StatusOK, responsePayload)
}

func InMemoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var record models.InMemoryRecord
		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, 8, "Error reading request body")
			return
		}
		if err = json.Unmarshal(body, &record); err != nil {
			respondWithError(w, http.StatusBadRequest, 9, "Invalid JSON data")
			return
		}

		// Store record using the key
		InMemoryDb[record.Key] = models.InMemoryRecord{
			Key:   record.Key,
			Value: record.Value,
		}
		respondWithJSON(w, http.StatusCreated, models.InMemoryResponsePayload{
			Code:    0,
			Msg:     "Success",
			Records: []models.InMemoryRecord{record},
		})

	case http.MethodGet:
		// Get the 'key' from the query string
		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			respondWithError(w, http.StatusBadRequest, 10, "URL Param 'key' is missing")
			return
		}

		key := keys[0]
		value, exists := InMemoryDb[key]
		if !exists {
			respondWithError(w, http.StatusNotFound, 11, "Record not found")
			return
		}

		record := models.InMemoryRecord{Key: key, Value: value.Value}
		respondWithJSON(w, http.StatusOK, models.InMemoryResponsePayload{Code: 0, Msg: "Success", Records: []models.InMemoryRecord{record}})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling the response"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}


func respondWithError(w http.ResponseWriter, statusCode int, code int, message string) {
	respondWithJSON(w, statusCode, models.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
