# Stage 1: Build the Go binary
FROM golang:1.23.6 AS builder

# Create and switch to a working directory
WORKDIR /app

# Copy go.mod and go.sum first for efficient caching
COPY go.mod ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 2: Create the final runtime image
FROM alpine:latest

WORKDIR /app

# Copy just the built binary plus any needed configs
COPY --from=builder /app/server .
COPY --from=builder /app/configs/ ./configs/

# Expose our server port
EXPOSE 8080

# Run the server binary
CMD ["./server"] 