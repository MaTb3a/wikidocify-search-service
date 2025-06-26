# Makefile for document-system
.PHONY: build run test clean dev docker-up docker-down docker-build docker-run docker-logs docker-logs-all docker-restart-app help

# Build the server binary
build:
	@echo "Building server..."
	@mkdir -p bin
	@go build -o bin/server ./cmd/server

# Run the server
run: build
	@echo "Starting server..."
	@./bin/server

# Run in development mode (with file watching if available)
dev:
	@echo "Starting in development mode..."
	@go run ./cmd/server

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin/

# Start Docker services (database only)
docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d

# Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down

# Build Docker image for the application
docker-build:
	@echo "Building Docker image..."
	@docker-compose build app

# Build and run full stack (database + application)
docker-run:
	@echo "Building and starting full stack..."
	@docker-compose up --build -d

# Show Docker container logs
docker-logs:
	@echo "Showing application logs..."
	@docker-compose logs -f app

# Show all container logs
docker-logs-all:
	@echo "Showing all container logs..."
	@docker-compose logs -f

# Restart just the application container
docker-restart-app:
	@echo "Restarting application container..."
	@docker-compose restart app

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Show help
help:
	@echo "Available commands:"
	@echo "  build           - Build the server binary"
	@echo "  run             - Build and run the server"
	@echo "  dev             - Run server in development mode"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-up       - Start Docker services (DB only)"
	@echo "  docker-down     - Stop Docker services"
	@echo "  docker-build    - Build Docker image for app"
	@echo "  docker-run      - Build and run full stack"
	@echo "  docker-logs     - Show application logs"
	@echo "  docker-logs-all - Show all container logs"
	@echo "  docker-restart-app - Restart application container"
	@echo "  deps            - Install and tidy dependencies"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter"
	@echo "  help            - Show this help message" 