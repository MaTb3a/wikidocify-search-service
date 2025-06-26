# WikiDocify Search Service

A high-performance microservice for searching documents in the WikiDocify project, powered by Elasticsearch and Kafka.  
This service keeps the search index in sync with your main Document Service and provides a RESTful API for fast, full-text search.

---

## Features

- **Full-Text Search**: Search documents by title and content.
- **Elasticsearch Integration**: Optimized indexing and querying.
- **Kafka Sync**: Real-time updates from the Document Service via Kafka.
- **RESTful API**: Clean endpoints for search and sync management.
- **Health Checks**: Built-in health endpoints for monitoring.
- **Docker Support**: Easy deployment with Docker Compose.
- **Kibana UI**: Visualize and debug your Elasticsearch data.

---

## Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop) and [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://golang.org/) (for local development)
- The Document Service running and accessible (see `DOC_SERVICE_URL`)

---

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd wikidocify-search-service
```

### 2. Configure Environment

Copy the example environment file and edit as needed:

```bash
cp .env.example .env
```

Edit `.env` with your settings:

```env
# Document Service URL (should be accessible from the search-service container)
DOC_SERVICE_URL=http://host.docker.internal:8081

# Elasticsearch URL
ELASTICSEARCH_URL=http://elasticsearch:9200

# Kafka Broker URL
KAFKA_BROKER=kafka:9092

# Elasticsearch index name
ELASTICSEARCH_INDEX=wikidocify_documents

# Enable automatic synchronization on startup
ENABLE_SYNC=true
```

### 3. Build and Run with Docker Compose

This will start **Elasticsearch**, **Kibana**, **Kafka**, **Zookeeper**, and the **Search Service**.

```bash
docker compose up --build
```

- **Elasticsearch**: [http://localhost:9200](http://localhost:9200)
- **Kibana**: [http://localhost:5601](http://localhost:5601)
- **Search Service API**: [http://localhost:8080](http://localhost:8080)

> **Note:** The Document Service must be running and accessible at the URL specified in `DOC_SERVICE_URL`.

### 4. Check Logs

```bash
docker compose logs -f search-service
```

### 5. Run Locally (for Development)

Start only Elasticsearch (and optionally Kafka/Zookeeper) with Docker Compose:

```bash
docker compose up -d elasticsearch kafka zookeeper
```

Then run the service locally:

```bash
go mod tidy
go run ./cmd/server
```

---

## API Endpoints

### Search Documents

- **Search in both title and content**
  ```
  GET /api/v1/search?query=your-search-term
  ```
- **Search only in titles**
  ```
  GET /api/v1/search?query=your-search-term&type=title
  ```
- **Search only in content**
  ```
  GET /api/v1/search?query=your-search-term&type=content
  ```
- **Paginated search**
  ```
  GET /api/v1/search?query=your-search-term&limit=20&offset=10
  ```
- **Filter by author**
  ```
  GET /api/v1/search?query=your-search-term&author=john-doe
  ```

### Sync Management

- **Trigger full sync (sync all documents from Document Service)**
  ```
  POST /api/v1/sync/full
  ```
- **Sync specific document**
  ```
  POST /api/v1/sync/document/{id}
  ```
- **Delete document from search index**
  ```
  DELETE /api/v1/sync/document/{id}
  ```
- **Check sync status**
  ```
  GET /api/v1/sync/status
  ```

### Health Check

- **Check service health**
  ```
  GET /health
  ```

---

## How It Works

1. **Document Service** emits document changes to Kafka.
2. **Search Service** consumes these events, updating the Elasticsearch index in real time.
3. **Search Service** exposes a REST API for searching and managing the index.
4. **Kibana** provides a UI for exploring and debugging Elasticsearch data.

---

## Troubleshooting

- **Elasticsearch or Kibana not available?**
  - Check Docker Compose logs: `docker compose logs elasticsearch kibana`
- **Search Service can't reach Document Service?**
  - Make sure `DOC_SERVICE_URL` is correct and accessible from inside the container.
- **Kafka errors?**
  - Ensure both Kafka and Zookeeper are running and healthy.

---

## Best Practices

- Use environment variables for all configuration.
- Keep your `.env` file out of version control.
- Monitor health endpoints for production readiness.
