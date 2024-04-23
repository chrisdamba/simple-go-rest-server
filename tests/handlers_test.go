package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert" // Import the testify library for assertions

	"github.com/chrisdamba/simple-go-rest-server"
)

func TestFetchFromMongo_Success(t *testing.T) {
	// ... 
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
