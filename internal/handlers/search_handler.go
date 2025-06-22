// internal/handlers/search_handler.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"wikidocify-search-service/internal/models"
	"wikidocify-search-service/internal/services"
)

type SearchHandler struct {
	searchService *services.SearchService
}

func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// Search handles search requests
// @Summary Search documents
// @Description Search documents by title and/or content
// @Tags search
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param type query string false "Search type: title, content, or all" default(all)
// @Param limit query int false "Number of results to return" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Param author query string false "Filter by author"
// @Success 200 {object} models.SearchResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	var req models.SearchRequest
	
	// Bind query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid search parameters",
			"details": err.Error(),
		})
		return
	}

	// Validate search query
	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query is required",
		})
		return
	}

	// Perform search
	result, err := h.searchService.Search(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Search failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SyncDocument handles manual document sync
// @Summary Sync a specific document
// @Description Manually sync a document from the doc service to Elasticsearch
// @Tags sync
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/sync/document/{id} [post]
func (h *SearchHandler) SyncDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid document ID",
		})
		return
	}

	if err := h.searchService.SyncDocument(uint32(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to sync document",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document synced successfully",
		"document_id": id,
	})
}

// DeleteDocument handles document deletion from search index
// @Summary Delete a document from search index
// @Description Remove a document from the Elasticsearch index
// @Tags sync
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/sync/document/{id} [delete]
func (h *SearchHandler) DeleteDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid document ID",
		})
		return
	}

	if err := h.searchService.DeleteDocument(uint32(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete document",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document deleted successfully",
		"document_id": id,
	})
}

// FullSync handles full synchronization
// @Summary Perform full sync
// @Description Sync all documents from the doc service to Elasticsearch
// @Tags sync
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/sync/full [post]
func (h *SearchHandler) FullSync(c *gin.Context) {
	if err := h.searchService.FullSync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Full sync failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Full sync completed successfully",
		"timestamp": time.Now(),
	})
}

// GetSyncStatus returns the current sync status
// @Summary Get sync status
// @Description Get information about sync configuration and last sync time
// @Tags sync
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Router /api/v1/sync/status [get]
func (h *SearchHandler) GetSyncStatus(c *gin.Context) {
	status := h.searchService.GetSyncStatus()
	c.JSON(http.StatusOK, gin.H{
		"sync_status": status,
	})
}

// Health handles health check requests
// @Summary Health check
// @Description Check the health of the search service and its dependencies
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Failure 503 {object} models.HealthResponse
// @Router /health [get]
func (h *SearchHandler) Health(c *gin.Context) {
	healthStatus := h.searchService.HealthCheck()
	
	response := models.HealthResponse{
		Timestamp:       time.Now(),
		ElasticsearchOK: healthStatus["elasticsearch"],
		DocServiceOK:    healthStatus["doc_service"],
		Details:         make(map[string]string),
	}

	// Determine overall status
	if response.ElasticsearchOK && response.DocServiceOK {
		response.Status = "healthy"
		c.JSON(http.StatusOK, response)
	} else {
		response.Status = "unhealthy"
		if !response.ElasticsearchOK {
			response.Details["elasticsearch"] = "Connection failed"
		}
		if !response.DocServiceOK {
			response.Details["doc_service"] = "Connection failed"
		}
		c.JSON(http.StatusServiceUnavailable, response)
	}
}