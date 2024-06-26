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


## Getting Started

To get started with this project, you need to have Go installed on your machine as well as MongoDB if you wish to run a local database server.

## Prerequisites

- Golang (version 1.19 or above)
- A running MongoDB instance


## Installation

Follow these steps to install the API on your local machine:

1. **Clone the repository**
   ```sh
   git clone https://github.com/chrisdamba/simple-go-rest-server.git
   cd simple-go-rest-server
   ```

2. **Set up the environment variable for MongoDB URI**
   - If using a local MongoDB instance, ensure it's running and note the URI connection string.
   - The database name should be `getir-case-study` and the collection name should be `records`.

   ```bash
   export MONGO_URI="mongodb://localhost:27017"
   ```

3. **Populate the database for testing**
   - You can use the MongoDB shell or a GUI like MongoDB Compass to insert test documents.
   ```javascript
   use getir-case-study
   db.records.insertMany([
     {
       "key": "TAKwGc6Jr4i8Z487",
       "createdAt": ISODate("2017-01-28T01:22:14.398Z"),
       "count": [500, 400, 450, 550, 300, 150, 350]
     },
     {
       "key": "NAeQ8eX7e5TEg70H",
       "createdAt": ISODate("2017-01-27T08:19:14.135Z"),
       "count": [540, 400, 450, 550, 300, 160, 350]
     },
     {
       "key": "cCddT2RPqWmUI4Nf",
       "createdAt": ISODate("2017-01-27T13:22:10.421Z"),
       "count": [120, 400, 450, 660, 500, 770, 250]
     }
   ])
   ```

4. **Build the project**
   ```bash
   go build -o api-server
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
curl --request POST \
  --url http://localhost:8080/mongo \
  --header 'Content-Type: application/json' \
  --data '{"startDate":"2016-01-26", "endDate":"2018-02-02", "minCount": 2700, "maxCount": 3000}'
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

mongodb+srv://challengeUser:WUMgIwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true