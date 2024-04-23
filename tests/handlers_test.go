package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/chrisdamba/simple-go-rest-server/handlers"
	"github.com/chrisdamba/simple-go-rest-server/models"
)

func TestFetchFromMongo_Success(t *testing.T) {
	// Mock request body
	requestBody, _ := json.Marshal(models.RequestPayload{
		StartDate: "2016-01-26",
		EndDate:   "2018-02-02",
		MinCount:  2700,
		MaxCount:  3000,
	})

	// Create a test request
	req, err := http.NewRequest("POST", "/mongo", bytes.NewBuffer(requestBody)) 
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture response
	rr := httptest.NewRecorder() 
	handler := http.HandlerFunc(handlers.FetchFromMongo)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Assert status code and response
	assert.Equal(t, http.StatusOK, rr.Code) 

	// Unmarshal response body
	var responsePayload models.ResponsePayload
	err = json.Unmarshal(rr.Body.Bytes(), &responsePayload)
	if err != nil {
		t.Fatal(err)
	}

	// More assertions on responsePayload.Code, responsePayload.Msg, responsePayload.Records
	assert.Equal(t, 0, responsePayload.Code)
	assert.Equal(t, "Success", responsePayload.Msg)
}

// Test cases for error scenarios, in-memory handler tests follow a similar pattern 
func TestFetchFromMongo_InvalidPayload(t *testing.T) {
	// ...
}

func TestInMemoryHandler_CreateRecord(t *testing.T) {
	// ...
}

func TestInMemoryHandler_GetRecords(t *testing.T) {
	// ...
}
