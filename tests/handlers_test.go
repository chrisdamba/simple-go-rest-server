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

// helper function to perform a POST request and return the response
func performRequest(handler http.HandlerFunc, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

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
	// Define an invalid payload (e.g., missing required fields)
	invalidPayload := map[string]string{
		"startDate": "2020-01-01", // Let's assume endDate is required and is missing
	}

	// Marshal the invalid payload into JSON
	payloadBytes, err := json.Marshal(invalidPayload)
	assert.NoError(t, err)

	// Create a new HTTP POST request with the invalid payload
	req, err := http.NewRequest("POST", "/mongo/fetch", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Instantiate a new router (assuming you have a function newRouter() that sets up your routes)
	handler := http.HandlerFunc(handlers.FetchFromMongo)

	// Call the ServeHTTP method directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect
	expectedErrorMessage := "Invalid request payload" // The error message your handler is expected to return
	assert.Contains(t, rr.Body.String(), expectedErrorMessage, "handler returned unexpected body")
}

func TestFetchFromMongo_MissingDates(t *testing.T) {
	handler := http.HandlerFunc(handlers.FetchFromMongo)

	// Missing StartDate and EndDate
	payload := map[string]int{"minCount": 2700, "maxCount": 3000}
	payloadBytes, _ := json.Marshal(payload)

	response := performRequest(handler, "POST", "/mongo", payloadBytes)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var resp map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &resp)

	assert.Equal(t, float64(2), resp["code"])
	assert.Equal(t, "StartDate and EndDate fields are required", resp["msg"])
}


func TestFetchFromMongo_InvalidStartDate(t *testing.T) {
	handler := http.HandlerFunc(handlers.FetchFromMongo) // your handler function here

	// Invalid StartDate format
	payload := models.RequestPayload{
		StartDate: "01-26-2016", // wrong format
		EndDate:   "2018-02-02",
		MinCount:  2700,
		MaxCount:  3000,
	}
	payloadBytes, _ := json.Marshal(payload)

	response := performRequest(handler, "POST", "/mongo", payloadBytes)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var resp map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &resp)

	assert.Equal(t, float64(3), resp["code"])
	assert.Equal(t, "Invalid StartDate format. Please use YYYY-MM-DD.", resp["msg"])
}

func TestFetchFromMongo_InvalidEndDate(t *testing.T) {
	handler := http.HandlerFunc(handlers.FetchFromMongo) // your handler function here

	// Invalid EndDate format
	payload := models.RequestPayload{
		StartDate: "2016-01-26",
		EndDate:   "02-02-2018", // wrong format
		MinCount:  2700,
		MaxCount:  3000,
	}
	payloadBytes, _ := json.Marshal(payload)

	response := performRequest(handler, "POST", "/mongo", payloadBytes)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var resp map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &resp)

	assert.Equal(t, float64(4), resp["code"])
	assert.Equal(t, "Invalid EndDate format. Please use YYYY-MM-DD.", resp["msg"])
}

func TestInMemoryHandler_CreateRecord(t *testing.T) {
	// ...
}

func TestInMemoryHandler_GetRecords(t *testing.T) {
	// ...
}
