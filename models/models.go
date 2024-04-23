package models

import (
	"time"
)

// MongoRecord represents a record in the MongoDB collection
type MongoRecord struct {
	Key        string    `bson:"key"`
	CreatedAt  time.Time `bson:"createdAt"`
	TotalCount int       `bson:"totalCount"`
	ID         string    `bson:"id"`
}

// InMemoryRecord represents a record in the in-memory database
type InMemoryRecord struct {
	ID          string `json:"id"`
	Data        string `json:"data"`
}

// RequestPayload for the MongoDB data fetch endpoint
type RequestPayload struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

// ResponsePayload is the format for the response of both endpoints
type ResponsePayload struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Records []struct {
		Key      string `json:"key,omitempty"`
		CreatedAt string `json:"createdAt,omitempty"`
		TotalCount int   `json:"totalCount,omitempty"`
	} `json:"records"`
}


type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}