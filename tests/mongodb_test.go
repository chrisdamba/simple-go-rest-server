package tests

import (
	"context"
	"testing"

	"github.com/chrisdamba/simple-go-rest-server/db"
	"github.com/stretchr/testify/assert"
)

func TestConnectToMongo_Success(t *testing.T) {
	// Connect to MongoDB
	mongoURI := "mongodb://localhost:27017"
	client, err := db.ConnectToMongo(mongoURI)
	assert.NoError(t, err)

	// Don't proceed if the connection failed
	if err != nil {
		assert.NoError(t, err)
		return
	}

	// Defer the disconnection of the client
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			t.Fatalf("Error on disconnecting from MongoDB: %v", err)
		}
	}()
	
	// Test if the client can ping the database
	err = client.Ping(context.TODO(), nil)
	assert.NoError(t, err)
}
