# items-packs-calculator

A simple HTTP-based service for calculating how best to fulfill an order quantity using predefined pack sizes. The project is written in Go and includes both a backend for handling pack calculations and a frontend for entering the desired item count.

## Features

1. **HTTP Endpoint:**  
   - `/calculate` (POST) – Receives a JSON payload like `{"items": 501}` and returns a JSON response indicating the pack distribution.

2. **Pack Sizes Configuration:**  
   - Stored in `configs/packs.json`  
   - Dynamically loaded at startup, with basic validations in place (e.g., pack sizes must be positive).

3. **Frontend:**  
   - A simple `index.html` that calls the local `/calculate` endpoint and displays results.

4. **Graceful Shutdown:**  
   - The server listens for OS signals (e.g., CTRL+C) and gracefully closes all resources.

## Project Structure

```
├── api/                      # Public HTTP API 
│   └── handlers.go
├── cmd/
│   └── server/
│       └── main.go           # Program entry point
├── configs/                  # JSON configs (packs.json)
├── internal/
│   ├── config/               # Configuration loading
│   └── packcalculator/       # Core pack calculation logic
├── frontend/                 # Simple HTML page to call API
├── Makefile                  # Basic build/test/run targets
├── Dockerfile                # Multi-stage Docker build
└── README.md                 # This file
```

## Getting Started

1. **Install Go**  
   - Requires Go 1.22 or newer (preferably 1.23.6).

2. **Clone & Build**
   ```bash
   git clone https://github.com/boriszhilko/items-packs-calculator.git
   cd items-packs-calculator
   make build
   ```

3. **Run**
   ```bash
   make run
   ```
   - The backend server will start at `http://localhost:8080`.

## Test

### Unit Tests

You can run all unit tests (excluding integration tests) with:

```bash
make test  # This will skip /test/ folder and run only unit tests
```

### Integration Tests

The integration tests spin up an in-memory server to test the full HTTP flow. To run them:

```bash
make integration-test
```

## Using the Application

### Via Frontend

1. Visit the deployed Heroku app at [https://items-packs-calculator-f97c28fbf434.herokuapp.com](https://items-packs-calculator-f97c28fbf434.herokuapp.com)
2. Enter a number of items (e.g., 501).
3. Click **Calculate**. A JSON result indicating your pack distribution will appear in the "Result" section.

### Via Curl

```bash
curl -X POST "https://items-packs-calculator-f97c28fbf434.herokuapp.com/calculate" \
  -H "Content-Type: application/json" \
  -d '{"items": 501}'
```

You'll receive a JSON response, for example:
```json
{
  "pack_distribution": {
    "500": 1,
    "250": 1
  },
  "total_items": 750
}
```

## Docker Instructions

1. **Build Docker Image**  
   ```bash
   docker build -t items-packs-calculator .
   ```
2. **Run Docker Container**  
   ```bash
   docker run -p 8080:8080 items-packs-calculator
   ```
3. **Access the API** at `http://localhost:8080/calculate`.