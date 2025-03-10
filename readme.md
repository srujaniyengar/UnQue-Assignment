# UnQue API

A simple appointment booking API built with Golang, Gin, and MongoDB.

## Setup

```bash
# Install dependencies
go mod download

# Ensure MongoDB is running at mongodb://localhost:27017

# Run the server
go run main.go
```
## Endpoints

    POST /login
    POST /availability
    GET /availability
    POST /appointments
    DELETE /appointments/:id
    GET /appointments


## Testing

```bash

# Run end-to-end tests
go test ./test -v
```
