# Document System

A Go-based document management system with PostgreSQL database and Docker support.

## Features

- Document upload and management
- PostgreSQL database integration
- Docker containerization
- RESTful API endpoints

## Prerequisites

- Go 1.19 or higher
- Docker and Docker Compose
- PostgreSQL

## Project Structure

```
document-system/
├── cmd/
│   └── server/         # Main application entry point
├── internal/
│   ├── config/        # Configuration files
│   ├── handlers/      # HTTP request handlers
│   ├── models/        # Database models
│   └── routes/        # API route definitions
├── docker-compose.yml # Docker configuration
└── Makefile          # Build and run commands
```

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/atheeralattar/pbl-week2.git
cd document-system
```

2. Start the application using Docker:
```bash
make docker-run
```

3. Check the logs to see if the application is running:
```bash
make docker-logs-all
```

The application will be available at `http://localhost:3976`

## Available Commands

- `make docker-run` - Build and run the full stack
- `make docker-logs` - View application logs
- `make docker-down` - Stop all containers
- `make docker-restart-app` - Restart just the application

## API Endpoints

- `POST /documents` - Upload a new document (title, content, author)
- `GET /documents` - List all documents
- `GET /documents/:id` - Get a specific document
- `PUT /documents/:id` - Update a document
- `DELETE /documents/:id` - Delete a document

## Development

To run the application in development mode:

```bash
make dev
```

For more commands, see the Makefile or run:
```bash
make help
```
