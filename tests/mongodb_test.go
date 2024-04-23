package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/chrisdamba/simple-go-rest-server"

)

func TestConnectToMongo_Success(t *testing.T) {
	client, err := main.connectToMongo()
	assert.NoError(t, err)

	// Test if the client can ping the database
	err = client.Ping(context.TODO(), nil)
	assert.NoError(t, err)
}
