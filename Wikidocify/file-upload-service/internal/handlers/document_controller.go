// controllers/document_controller.go
// This file contains the controller for the document resource
// It defines the DocumentController struct and the NewDocumentController function
// It also contains the methods for the DocumentController struct
// It uses the Gin framework to handle the HTTP requests
// It uses the GORM library to interact with the database

package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"wikidocify/file-upload-service/internal/kafka"
	"wikidocify/file-upload-service/internal/models"

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

	document := models.Document{
		Title:   docRequest.Title,
		Content: []byte(docRequest.Content),
		Author:  docRequest.Author,
	}

	log.Printf("[API] POST /documents - Received document: title='%s', author='%s', content_size=%d bytes",
		document.Title, document.Author, len(document.Content))

	if err := dc.documentModel.Create(&document); err != nil {
		log.Printf("[DATABASE] Failed to create document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[DATABASE] Document created successfully with ID: %d", document.ID)

	// Publish Kafka event for document creation
	if err := kafka.PublishDocEvent("created", strconv.FormatUint(uint64(document.ID), 10), document.Title, string(document.Content)); err != nil {
		log.Printf("[KAFKA] Failed to publish created event for document ID %d: %v", document.ID, err)
	} else {
		log.Printf("[KAFKA] Successfully published created event for document ID %d", document.ID)
	}

	log.Printf("[API] POST /documents - Returning created document (ID: %d)", document.ID)
	c.JSON(http.StatusCreated, document)
}

func (dc *DocumentController) GetAll(c *gin.Context) {
	log.Printf("[API] GET /documents - Fetching all documents from IP: %s", c.ClientIP())

	// Parse pagination params
	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	log.Printf("[DATABASE] Querying documents: page=%d, limit=%d", page, limit)
	documents, total, err := dc.documentModel.FindAllPaginated(page, limit)
	if err != nil {
		log.Printf("[DATABASE] Failed to fetch documents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
		"page":      page,
		"limit":     limit,
		"total":     total,
	})
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

	var docRequest DocumentRequest
	if err := c.ShouldBindJSON(&docRequest); err != nil {
		log.Printf("[API] PUT /documents/%s - Invalid JSON payload: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	document.Title = docRequest.Title
	document.Content = []byte(docRequest.Content)
	document.Author = docRequest.Author

	log.Printf("[API] PUT /documents/%s - New data: title='%s', author='%s', content_size=%d bytes",
		id, document.Title, document.Author, len(document.Content))

	if err := dc.documentModel.Update(&document); err != nil {
		log.Printf("[DATABASE] Failed to update document with ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[DATABASE] Document with ID %s updated successfully", id)

	// Publish Kafka event for document update
	if err := kafka.PublishDocEvent("updated", strconv.FormatUint(uint64(document.ID), 10), document.Title, string(document.Content)); err != nil {
		log.Printf("[KAFKA] Failed to publish updated event for document ID %s: %v", id, err)
	} else {
		log.Printf("[KAFKA] Successfully published updated event for document ID %s", id)
	}

	log.Printf("[API] PUT /documents/%s - Returning updated document", id)
	c.JSON(http.StatusOK, document)
}

func (dc *DocumentController) Delete(c *gin.Context) {
	id := c.Param("id")
	log.Printf("[API] DELETE /documents/%s - Deleting document from IP: %s", id, c.ClientIP())

	document, err := dc.documentModel.FindByID(id)
	if err != nil {
		log.Printf("[DATABASE] Document with ID %s not found for deletion: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	log.Printf("[DATABASE] Document with ID %s found, proceeding with deletion...", id)

	if err := dc.documentModel.Delete(id); err != nil {
		log.Printf("[DATABASE] Failed to delete document with ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[DATABASE] Document with ID %s deleted successfully", id)

	// Publish Kafka event for document deletion
	if err := kafka.PublishDocEvent("deleted", id, document.Title, string(document.Content)); err != nil {
		log.Printf("[KAFKA] Failed to publish deleted event for document ID %s: %v", id, err)
	} else {
		log.Printf("[KAFKA] Successfully published deleted event for document ID %s", id)
	}

	log.Printf("[API] DELETE /documents/%s - Document deleted successfully", id)
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted"})
}
