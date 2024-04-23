# Simple Go REST Server

## Introduction

This project implements a simple RESTful API in Go demonstrating MongoDB integration and in-memory data storage without using any external web frameworks. It includes two primary functionalities:

1. Fetching data from a MongoDB database based on query parameters passed through a POST request.
2. An in-memory datastore that can be accessed and modified via GET and POST requests.

This API is designed to demonstrate the ability to create RESTful services in Go using only the standard library and a MongoDB driver for database interactions.



## Features

- **MongoDB Integration:**
  - Fetches data from a MongoDB collection based on provided criteria (date range, counts).
- **In-Memory Data Store:**
  - Creates in-memory records with generated IDs.
  - Fetches all stored in-memory records.
- **REST API Endpoints:**
  - `/mongo` (POST) - Fetches from MongoDB.
  - `/in-memory` (POST) - Creates an in-memory record.
  - `/in-memory` (GET) - Gets all in-memory records.


## Getting Started

To get started with this project, you need to have Go installed on your machine as well as MongoDB if you wish to run a local database server.

## Prerequisites

- Golang (version 1.19 or above)
- A running MongoDB instance


## Installation

Follow these steps to install the API on your local machine:

1.  **Clone the repository**

```sh
git clone https://github.com/chrisdamba/simple-go-rest-server.git
cd simple-go-rest-server
```

2.  **Set up the environment variable for MongoDB URI**
```bash
export MONGO_URI="your_mongodb_uri"
```

3.  **Build the project**
```bash
go build -o 
```



## Running the Application

After installation, you can run the application with the following command:

```sh
./api-server
```

By default, the API server runs on `http://localhost:8080`. This can be configured in the source code if necessary.

## Endpoints

The API has the following endpoints:

- POST `/mongo`
  - Fetches data from the MongoDB database. Requires a JSON body with `startDate`, `endDate`, `minCount`, and `maxCount`.

- POST `/in-memory`
  - Stores a key-value pair in the in-memory datastore. Requires a JSON body with `key` and `value`.

- GET `/in-memory`
  - Retrieves a value by key from the in-memory datastore. Requires a query parameter `key`.

## Testing

To test the endpoints, you can use tools like `curl` or Postman.

Example `curl` commands:

Fetch from MongoDB:
```sh
curl -X POST http://localhost:8080/mongo -d '{"startDate":"2016-01-26", "endDate":"2018-02-02", "minCount": 2700, "maxCount": 3000}' -H "Content-Type: application/json"
```

Add to in-memory datastore:
```sh
curl -X POST http://localhost:8080/in-memory -d '{"key":"exampleKey", "value":"exampleValue"}' -H "Content-Type: application/json"
```

Retrieve from in-memory datastore:
```sh
curl http://localhost:8080/in-memory?key=exampleKey
```

## Conclusion

This project showcases the capability of Go to build RESTful APIs without relying on third-party frameworks, using its powerful standard library and efficient concurrency model.
