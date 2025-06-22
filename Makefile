# Makefile
.PHONY: help build run test clean docker-build docker-run docker-stop logs

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the Go binary"
	@echo "  run          - Run the service locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker containers"
	@echo "  logs         - Show Docker logs"
	@echo "  sync-full    - Trigger full sync"
	@echo "  health       - Check service health"

# Build the Go binary
build:
	@echo "Building search service..."
	go build -o bin/search-service ./cmd/server

# Run the service locally
run:
	@echo "Starting search service..."
	go run ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker compose build

# Run with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	docker compose up -d

# Stop Docker containers
docker-stop:
	@echo "Stopping Docker containers..."
	docker compose down

# Show Docker logs
logs:
	@echo "Showing logs..."
	docker compose logs -f search-service

# Trigger full sync
sync-full:
	@echo "Triggering full sync..."
	curl -X POST http://localhost:8080/api/v1/sync/full

# Check service health
health:
	@echo "Checking service health..."
	curl -X GET http://localhost:8080/health | jq

# Search example
search-example:
	@echo "Example search..."
	curl "http://localhost:8080/api/v1/search?query=example&type=all&limit=5" | jq

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go mod tidy
	cp .env.example .env
	@echo "Please update .env with your configuration"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy