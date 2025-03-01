# Makefile for items-packs-calculator

# Binaries
BINARY_NAME = server
BIN_DIR = bin

.PHONY: all build test run clean

all: build test

build:
	@echo "==> Building the Go application..."
	@go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/server

test:
	@echo "==> Running tests..."
	@go test ./... -v

run: build
	@echo "==> Running the server..."
	@./$(BIN_DIR)/$(BINARY_NAME)

clean:
	@echo "==> Cleaning up build artifacts..."
	@rm -rf $(BIN_DIR) 