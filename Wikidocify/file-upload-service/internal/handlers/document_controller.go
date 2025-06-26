// controllers/document_controller.go
// This file contains the controller for the document resource
// It defines the DocumentController struct and the NewDocumentController function
// It also contains the methods for the DocumentController struct
// It uses the Gin framework to handle the HTTP requests
// It uses the GORM library to interact with the database

package handlers

import (
	"log"
	"net/http"

	"github.com/atheeralattar/pbl-week2/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentController struct {
	documentModel *models.DocumentModel
}

// DocumentRequest represents the JSON structure for incoming requests
type DocumentRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author"`
}

func NewDocumentController(db *gorm.DB) *DocumentController {
	log.Println("[CONTROLLER] Creating new DocumentController instance")
	return &DocumentController{
		documentModel: models.NewDocumentModel(db),
	}
}

func (dc *DocumentController) Create(c *gin.Context) {
	log.Printf("[API] POST /documents - Creating new document from IP: %s", c.ClientIP())

	var docRequest DocumentRequest
	if err := c.ShouldBindJSON(&docRequest); err != nil {
		log.Printf("[API] POST /documents - Invalid JSON payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert string content to []byte
	document := models.Document{
		Title:   docRequest.Title,
		Content: []byte(docRequest.Content),
		Author:  docRequest.Author,
	}

	log.Printf("[API] POST /documents - Received document: title='%s', author='%s', content_size=%d bytes",
		document.Title, document.Author, len(document.Content))

	// Save the document to the database
	log.Println("[DATABASE] Attempting to create document in database...")
	if err := dc.documentModel.Create(&document); err != nil {
		log.Printf("[DATABASE] Failed to create document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[DATABASE] Document created successfully with ID: %d", document.ID)
	log.Printf("[API] POST /documents - Returning created document (ID: %d)", document.ID)

	// Respond with the newly created document
	c.JSON(http.StatusCreated, document)
}

func (dc *DocumentController) GetAll(c *gin.Context) {
	log.Printf("[API] GET /documents - Fetching all documents from IP: %s", c.ClientIP())

	log.Println("[DATABASE] Querying all documents from database...")
	documents, err := dc.documentModel.FindAll()
	if err != nil {
		log.Printf("[DATABASE] Failed to fetch documents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[DATABASE] Successfully retrieved %d documents", len(documents))
	log.Printf("[API] GET /documents - Returning %d documents", len(documents))
	c.JSON(http.StatusOK, documents)
}

func (dc *DocumentController) GetByID(c *gin.Context) {
	id := c.Param("id")
	log.Printf("[API] GET /documents/%s - Fetching document by ID from IP: %s", id, c.ClientIP())

	log.Printf("[DATABASE] Querying document with ID: %s", id)
	document, err := dc.documentModel.FindByID(id)
	if err != nil {
		log.Printf("[DATABASE] Document with ID %s not found: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	log.Printf("[DATABASE] Found document: ID=%d, title='%s', content_size=%d bytes",
		document.ID, document.Title, len(document.Content))
	log.Printf("[API] GET /documents/%s - Returning document", id)
	c.JSON(http.StatusOK, document)
}

func (dc *DocumentController) Update(c *gin.Context) {
	id := c.Param("id")
	log.Printf("[API] PUT /documents/%s - Updating document from IP: %s", id, c.ClientIP())

	log.Printf("[DATABASE] Finding document with ID: %s for update", id)
	document, err := dc.documentModel.FindByID(id)
	if err != nil {
		log.Printf("[DATABASE] Document with ID %s not found for update: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	log.Printf("[DATABASE] Found document to update: ID=%d, current title='%s', content_size=%d bytes",
		document.ID, document.Title, len(document.Content))

	// Bind the incoming data
	var docRequest DocumentRequest
	if err := c.ShouldBindJSON(&docRequest); err != nil {
		log.Printf("[API] PUT /documents/%s - Invalid JSON payload: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update document fields
	document.Title = docRequest.Title
	document.Content = []byte(docRequest.Content)
	document.Author = docRequest.Author

	log.Printf("[API] PUT /documents/%s - New data: title='%s', author='%s', content_size=%d bytes",
		id, document.Title, document.Author, len(document.Content))

	// Update the document in the database
	log.Printf("[DATABASE] Updating document with ID: %s in database...", id)
	if err := dc.documentModel.Update(&document); err != nil {
		log.Printf("[DATABASE] Failed to update document with ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[DATABASE] Document with ID %s updated successfully", id)
	log.Printf("[API] PUT /documents/%s - Returning updated document", id)
	c.JSON(http.StatusOK, document)
}

func (dc *DocumentController) Delete(c *gin.Context) {
	id := c.Param("id")
	log.Printf("[API] DELETE /documents/%s - Deleting document from IP: %s", id, c.ClientIP())

	// Check if document exists
	log.Printf("[DATABASE] Checking if document with ID %s exists...", id)
	_, err := dc.documentModel.FindByID(id)
	if err != nil {
		log.Printf("[DATABASE] Document with ID %s not found for deletion: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	log.Printf("[DATABASE] Document with ID %s found, proceeding with deletion...", id)

	// Delete the document
	if err := dc.documentModel.Delete(id); err != nil {
		log.Printf("[DATABASE] Failed to delete document with ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[DATABASE] Document with ID %s deleted successfully", id)
	log.Printf("[API] DELETE /documents/%s - Document deleted successfully", id)
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted"})
}
