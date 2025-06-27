// internal/routes/routes.go
package routes

import (
	"wikidocify/elasticsearch-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all HTTP routes for the search service.
func SetupRoutes(router *gin.Engine, searchHandler *handlers.SearchHandler) {
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check (no API prefix)
	router.GET("/health", searchHandler.Health)

	// Root endpoint for service info
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "wikidocify-search-service",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	api := router.Group("/api/v1")
	{
		api.GET("/search", searchHandler.Search)
	}
}