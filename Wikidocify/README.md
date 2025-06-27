# WikiDocify Microservices Platform

A modern, production-ready document management and search platform built with Go, PostgreSQL, Kafka, and Elasticsearch.  
This project demonstrates best practices for microservices architecture, event-driven data sync, and scalable search.

---

## Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Microservices](#microservices)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Environment Variables](#environment-variables)
  - [Running with Docker Compose](#running-with-docker-compose)
  - [Development Workflow](#development-workflow)
- [API Reference](#api-reference)
  - [File Upload Service](#file-upload-service-api)
  - [Search Service](#search-service-api)
- [Event Model](#event-model)
- [Health Checks](#health-checks)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

---

## Project Overview

**WikiDocify** is a microservices-based platform for document management and full-text search.  
It is designed for reliability, scalability, and extensibility, using event-driven architecture and modern DevOps practices.

**Key goals:**
- Decouple document storage from search/indexing.
- Enable real-time sync between services using Kafka.
- Provide robust REST APIs for document CRUD and search.
- Support horizontal scaling and easy observability.

---

## Architecture

GitHub Copilot
+---------------------+ Kafka +---------------------+ Elasticsearch +-------------------+ | File Upload Service | --------------> | Search Service | --------------------> | Elasticsearch | | (REST + DB + Kafka) | (events) | (REST + Kafka) | (indexing/search) | | +---------------------+ +---------------------+ +-------------------+ | ^ | | v | PostgreSQL | | | +-------------------<-------------------<-------------------<-------------------------+

- **File Upload Service**: Handles document CRUD, persists to PostgreSQL, emits events to Kafka.
- **Search Service**: Consumes Kafka events, syncs data to Elasticsearch, exposes search API.
- **Kafka**: Event bus for reliable, decoupled communication.
- **Elasticsearch**: High-performance search and analytics engine.
- **PostgreSQL**: Source of truth for document data.

---

## Microservices

### 1. **File Upload Service**
- RESTful API for document CRUD.
- Stores documents in PostgreSQL.
- Publishes document events (create/update/delete) to Kafka.

### 2. **Search Service**
- Consumes document events from Kafka.
- Syncs and indexes documents in Elasticsearch.
- Exposes REST API for full-text search and sync management.

---

## Features

- **Microservices architecture** with clear separation of concerns.
- **Event-driven sync** using Kafka for real-time updates.
- **Full-text search** with Elasticsearch.
- **Health checks** for all services.
- **Pagination** and filtering for large datasets.
- **Docker Compose** for easy local development and orchestration.
- **Kibana** and **Kafka UI** for observability and debugging.

---

## Tech Stack

- **Go** (Golang) for all backend services
- **PostgreSQL** for document storage
- **Kafka** for event streaming
- **Elasticsearch** for search
- **Docker & Docker Compose** for orchestration
- **Gin** for REST APIs
- **Kibana** for Elasticsearch UI
- **Kafka UI** for Kafka management

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://golang.org/) (for local development)

### Environment Variables

All configuration is managed via `.env` in the project root.  
See `.env` for all options and sensible defaults.

**Example:**
```env
# Database
DB_HOST=postgres
DB_USER=doc_db_admin
DB_PASSWORD=SecurePass889
DB_NAME=documents_db
DB_PORT=5432
DB_SSLMODE=disable

# Service Ports
UPLOAD_SERVICE_PORT=8081
SEARCH_SERVICE_PORT=8080

# Kafka
KAFKA_BROKER=kafka:9092
KAFKA_TOPIC=document-events

# Elasticsearch
ELASTICSEARCH_URL=http://elasticsearch:9200
ELASTICSEARCH_INDEX=wikidocify_documents

# Document Service (for search service)
DOC_SERVICE_URL=http://file-upload-service:8081
```


# create the topic manually
docker exec -it wikidocify-kafka kafka-topics --create --topic document-events --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

Or use Kafka UI at http://localhost:8090 to create document-events.

