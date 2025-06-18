# syntax=docker/dockerfile:1
FROM golang:1.21-alpine

# Install necessary tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the Go app
RUN go build -o search-service ./cmd/server

# Expose the service port
EXPOSE 8081

# Run the binary
CMD ["./search-service"]
