# Oblivion

Oblivion is a Go-based playground project that implements a REST API for managing buckets and their associated keys and values. The project demonstrates how to build a CRUD-based REST API using only Go's standard libraries and SQLite for data persistence.

## Features

- **Bucket Management**: Create, retrieve, list, and delete buckets.
- **Key-Value Management**: Store, retrieve, update, and delete key-value pairs within buckets.
- **Querying**: Search for keys based on criteria.
- **Validation**: Input validation for bucket names, field names, and data types.
- **Error Handling**: Structured error responses with HTTP status codes and error descriptions.

## Project Structure

The project is organized into the following directories:

- `api/`: Contains the REST API handlers and routing logic.
- `bucket/`: Implements the business logic for bucket operations.
- `httprouter/`: Provides a lightweight HTTP router and response utilities.
- `model/`: Defines the data models used in the application.
- `repo/`: Implements the repository layer for data persistence using SQLite.
- `valid/`: Contains validation logic for input data.
- `apperror/`: Defines application-specific error types and handling.

## REST API Endpoints

### Buckets

#### List Buckets
**GET** `/v1/buckets`

Response:
```json
[
  "bucket1",
  "bucket2"
]
```

#### Create Bucket
**POST** `/v1/buckets`

Request Body:
```json
{
  "name": "people",
  "schema": [
    {
      "field": "id",
      "type": "string",
      "not-null": true,
      "indexed": true
    },
    {
      "field": "first_name",
      "type": "string",
      "not-null": true,
      "indexed": false
    }
  ]
}
```

Response:
```json
{
  "name": "people",
  "schema": [
    {
      "field": "id",
      "type": "string",
      "not-null": true,
      "indexed": true
    },
    {
      "field": "first_name",
      "type": "string",
      "not-null": true,
      "indexed": false
    }
  ]
}
```

#### Get Bucket
**GET** `/v1/buckets/{bucket}`

Response:
```json
{
  "name": "people",
  "schema": [
    {
      "field": "id",
      "type": "string",
      "not-null": true,
      "indexed": true
    },
    {
      "field": "first_name",
      "type": "string",
      "not-null": true,
      "indexed": false
    }
  ]
}
```

#### Delete Bucket
**DELETE** `/v1/buckets/{bucket}`

Response: `204 No Content`

---

### Keys

#### Get Key
**GET** `/v1/buckets/{bucket}/keys/{key}`

Response:
```json
{
  "id": "id1",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Set Key
**PUT** `/v1/buckets/{bucket}/keys/{key}`

Request Body:
```json
{
  "id": "id1",
  "first_name": "John",
  "last_name": "Doe"
}
```

Response: `204 No Content`

#### Delete Key
**DELETE** `/v1/buckets/{bucket}/keys/{key}`

Response: `204 No Content`

#### Find Keys
**GET** `/v1/buckets/{bucket}/keys?field=value`

Response:
```json
[
  "key1",
  "key2"
]
```

## Running the Project

1. Install Go (version 1.22 or later).
2. Clone the repository:
   ```sh
   git clone https://github.com/jjmrocha/oblivion.git
   cd oblivion
   ```
3. Run the application:
   ```sh
   go run main.go
   ```
4. The server will start on `http://localhost:9090`.

## Testing the API

You can use the provided `test.http` file to test the API using tools like [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) in Visual Studio Code.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.