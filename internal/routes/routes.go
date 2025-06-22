// internal/routes/routes.go
package routes

import (
	"wikidocify-search-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, searchHandler *handlers.SearchHandler) {
	// Basic middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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
		// Only search endpoint
		search := api.Group("/search")
		{
			search.GET("", searchHandler.Search)
		}
	}
}