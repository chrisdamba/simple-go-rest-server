package models

import (
	"time"
)

// MongoRecord represents a record in the MongoDB collection
type MongoRecord struct {
	Key        string    `bson:"key"`
	CreatedAt  time.Time `bson:"createdAt"`
	Count 		 []int     `bson:"count"`
	ID         string    `bson:"id"`
}

// InMemoryRecord represents a record in the in-memory database
type InMemoryRecord struct {
	Key   	string `json:"key"`
	Value 	string `json:"value"`
}

type InMemoryResponsePayload struct {
	Code    int             `json:"code"`
	Msg     string          `json:"msg"`
	Records []InMemoryRecord `json:"records,omitempty"` // The records field is omitted if empty
}
// RequestPayload for the MongoDB data fetch endpoint
type RequestPayload struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

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