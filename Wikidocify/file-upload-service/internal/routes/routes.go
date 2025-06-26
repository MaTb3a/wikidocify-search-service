// routes/routes.go
// This file contains the routes for the document resource
// It defines the SetupRoutes function
// It uses the Gin framework to handle the HTTP requests
// It uses the GORM library to interact with the database

package routes

import (
	"log"

	"github.com/atheeralattar/pbl-week2/internal/config"
	"github.com/atheeralattar/pbl-week2/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	log.Println("[ROUTES] Initializing Gin router...")
	r := gin.Default()

	// Initialize controllers
	log.Println("[ROUTES] Creating document controller...")
	documentController := handlers.NewDocumentController(config.GetDB())

	// Document routes
	log.Println("[ROUTES] Setting up document endpoints...")
	documentRoutes := r.Group("/documents")
	{
		documentRoutes.POST("", documentController.Create)
		documentRoutes.GET("", documentController.GetAll)
		documentRoutes.GET("/:id", documentController.GetByID)
		documentRoutes.PUT("/:id", documentController.Update)
		documentRoutes.DELETE("/:id", documentController.Delete)
	}

	log.Println("[ROUTES] All routes registered successfully")
	return r
}
