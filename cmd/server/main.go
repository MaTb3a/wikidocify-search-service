// cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wikidocify-search-service/internal/config"
	"wikidocify-search-service/internal/elastic"
	"wikidocify-search-service/internal/handlers"
	"wikidocify-search-service/internal/routes"
	"wikidocify-search-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize Elasticsearch client
	esClient, err := elastic.NewClient(
		cfg.Elasticsearch.URL,
		cfg.Elasticsearch.Username,
		cfg.Elasticsearch.Password,
		cfg.Elasticsearch.Index,
	)
	if err != nil {
		log.Fatal("Failed to create Elasticsearch client:", err)
	}
	log.Println("Elasticsearch client initialized")

	// Initialize document service client
	docServiceClient := services.NewDocServiceClient(
		cfg.DocService.BaseURL,
		cfg.DocService.APIKey,
		cfg.DocService.Timeout,
	)
	log.Println("Document service client initialized")

	// Initialize search service
	searchService := services.NewSearchService(
		esClient,
		docServiceClient,
		cfg.Sync.SyncInterval,
		cfg.Sync.BatchSize,
		cfg.Sync.EnableSync,
	)
	log.Println("Search service initialized")

	// Perform initial full sync if enabled
	if cfg.Sync.EnableSync {
		log.Println("Starting initial full sync...")
		if err := searchService.FullSync(); err != nil {
			log.Printf("  Initial sync failed: %v", err)
		} else {
			log.Println(" Initial sync completed")
		}

		// Start periodic sync
		searchService.StartPeriodicSync()
		log.Println(" Periodic sync started")
	}

	// Initialize handlers
	searchHandler := handlers.NewSearchHandler(searchService)

	// Setup Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.New()
	routes.SetupRoutes(router, searchHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf(" Server starting on port %s", cfg.Server.Port)
		log.Printf(" Elasticsearch: %s", cfg.Elasticsearch.URL)
		log.Printf(" Document Service: %s", cfg.DocService.BaseURL)
		log.Printf(" Sync enabled: %v", cfg.Sync.EnableSync)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println(" Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println(" Server exited")
}