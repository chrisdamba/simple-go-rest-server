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

	assert.Equal(t, 0, responsePayload.Code)
	assert.Equal(t, "Success", responsePayload.Msg)
}

// Test cases for error scenarios
func TestFetchFromMongo_InvalidPayload(t *testing.T) {
	handler := http.HandlerFunc(handlers.FetchFromMongo) 
	// Define an invalid payload (e.g., missing required fields)
	invalidPayload := map[string]string{
		"startDate": "2020-01-01", // assume endDate is required and is missing
	}

	// Marshal the invalid payload into JSON
	payloadBytes, err := json.Marshal(invalidPayload)
	assert.NoError(t, err)

	rr := performRequest(handler, "POST", "/mongo", payloadBytes)
	// Check the status code is what we expect
	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect
	expectedErrorMessage := "Invalid request payload" 
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
	handler := http.HandlerFunc(handlers.FetchFromMongo) 

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
	handler := http.HandlerFunc(handlers.FetchFromMongo) 

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
	// Prepare the server and handler
	server := httptest.NewServer(http.HandlerFunc(handlers.InMemoryHandler))
	defer server.Close()

	// Create a sample record to POST
	record := models.InMemoryRecord{Key: "testKey", Value: "testValue"}
	recordJSON, _ := json.Marshal(record)

	// Create a request to our handler
	response, err := http.Post(server.URL, "application/json", bytes.NewBuffer(recordJSON))
	assert.NoError(t, err)

	// Check the status code
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	// Decode the response body
	var respPayload models.InMemoryResponsePayload
	err = json.NewDecoder(response.Body).Decode(&respPayload)
	assert.NoError(t, err)

	// Validate the response payload
	assert.Equal(t, 0, respPayload.Code)
	assert.Equal(t, "Success", respPayload.Msg)
	assert.NotEmpty(t, respPayload.Records)
	assert.Equal(t, "testKey", respPayload.Records[0].Key)
	assert.Equal(t, "testValue", respPayload.Records[0].Value)
}

func TestInMemoryHandler_GetRecords(t *testing.T) {
	// Prepare the server and handler
	server := httptest.NewServer(http.HandlerFunc(handlers.InMemoryHandler))
	defer server.Close()

	// Assume the inMemoryDb is already populated with a record
	handlers.InMemoryDb["testKey"] = models.InMemoryRecord{Key: "testKey", Value: "testValue"}

	// Build the GET request
	req, _ := http.NewRequest("GET", server.URL+"?key=testKey", nil)

	// Create a response recorder (to record the response)
	responseRecorder := httptest.NewRecorder()

	// Dispatch the request to our handler
	handler := http.HandlerFunc(handlers.InMemoryHandler)
	handler.ServeHTTP(responseRecorder, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// Decode the response body
	var respPayload models.InMemoryResponsePayload
	err := json.NewDecoder(responseRecorder.Body).Decode(&respPayload)
	assert.NoError(t, err)

	// Validate the response payload
	assert.Equal(t, 0, respPayload.Code)
	assert.Equal(t, "Success", respPayload.Msg)
	assert.NotEmpty(t, respPayload.Records)
	assert.Equal(t, "testKey", respPayload.Records[0].Key)
	assert.Equal(t, "testValue", respPayload.Records[0].Value)
}