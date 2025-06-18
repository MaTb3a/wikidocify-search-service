package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"wikidocify-search-service/internal/config"
	"wikidocify-search-service/internal/elastic"
	"wikidocify-search-service/internal/handlers"
	"wikidocify-search-service/internal/kafka"
	"wikidocify-search-service/internal/routes"
)

func main() {
	log.Println("🟢 Starting Search Service...")

	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found. Continuing with system environment...")
	}
	config.LoadEnv()

	// Init Elastic
	elastic.InitElasticClient()

	// Start Kafka listener (async)
	go kafka.StartConsumer()

	// Setup server
	router := gin.Default()
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Println("🚀 Running on port", port)
	router.Run(":" + port)
}
