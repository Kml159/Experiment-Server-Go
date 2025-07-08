# Experiment Server Go

A Go server for distributing experiment parameters to clients and collecting results.

## Features
- Client management and status tracking
- Experiment parameter generation and distribution
- Result collection from completed experiments

## Structure
```
├── cmd/server/main.go           # Server entry point
├── internal/models/
│   ├── client/client.go         # Client model
│   └── parameter/parameter.go   # Parameter generation
```

## Models
- **Client**: Tracks experiment status, completion history, and connection info
- **Parameter**: Defines genetic algorithm parameters for optimization experiments

## Running the Server

```sh
go run cmd/server/main.go
```

## Building
```sh
go build -o experiment-server cmd/server/main.go
```

## Requirements
- Go 1.18+

