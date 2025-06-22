// internal/routes/routes.go
package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"wikidocify-search-service/internal/handlers"
	"wikidocify-search-service/internal/middleware"
)

func SetupRoutes(router *gin.Engine, searchHandler *handlers.SearchHandler) {
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())
	router.Use(middleware.RateLimit(100, time.Minute)) // 100 requests per minute

	// Health check (no API prefix)
	router.GET("/health", searchHandler.Health)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "wikidocify-search-service",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Search endpoints
		search := api.Group("/search")
		{
			search.GET("", searchHandler.Search)
		}

		// Sync endpoints (protected by API key if needed)
		sync := api.Group("/sync")
		sync.Use(middleware.APIKeyAuth()) // Optional: protect sync endpoints
		{
			sync.POST("/document/:id", searchHandler.SyncDocument)
			sync.DELETE("/document/:id", searchHandler.DeleteDocument)
			sync.POST("/full", searchHandler.FullSync)
			sync.GET("/status", searchHandler.GetSyncStatus)
		}
	}
}