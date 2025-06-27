package main

import (
	"log"
	"os"

	"wikidocify/file-upload-service/internal/config"
	"wikidocify/file-upload-service/internal/kafka"
	"wikidocify/file-upload-service/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("[STARTUP] Starting file-upload-service server...")

	// Load environment variables from .env file
	log.Println("[CONFIG] Loading environment variables...")
	if err := godotenv.Load(); err != nil {
		log.Printf("[CONFIG] Warning: Error loading .env file: %v", err)
		log.Println("[CONFIG] Continuing with system environment variables...")
	} else {
		log.Println("[CONFIG] Environment variables loaded successfully")
	}

	// Initialize database
	log.Println("[DATABASE] Initializing database connection...")
	config.InitDB()

	// Initialize Kafka producer
	log.Println("[KAFKA] Initializing Kafka producer...")
	kafka.InitKafkaWriter()

	// Setup routes (CORS is handled inside SetupRoutes)
	log.Println("[ROUTES] Setting up HTTP routes...")
	router := routes.SetupRoutes()

	// Get server port from environment variable
	port := os.Getenv("UPLOAD_SERVICE_PORT")
	if port == "" {
		port = "8081" // fallback to default port
		log.Printf("[CONFIG] UPLOAD_SERVICE_PORT not set, using default port: %s", port)
	} else {
		log.Printf("[CONFIG] Using configured port: %s", port)
	}

	// Start the server
	log.Printf("[SERVER] Starting HTTP server on port %s...", port)
	log.Printf("[SERVER] Server ready! Visit http://localhost:%s/documents", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("[SERVER] Failed to start server: %v", err)
	}
}

// ENV GOPROXY=direct
