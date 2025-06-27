package handlers

import (
	"net/http"
	"time"

	"wikidocify/elasticsearch-service/internal/models"
	"wikidocify/elasticsearch-service/internal/services"

	"github.com/gin-gonic/gin"
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
func (h *SearchHandler) Search(c *gin.Context) {
    var req models.SearchRequest

    // Bind query parameters
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid search parameters",
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

    // Set default pagination if not provided
    if req.Limit == 0 {
        req.Limit = 10
    }
    if req.Offset < 0 {
        req.Offset = 0
    }

    // Perform search
    docs, total, err := h.searchService.Search(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Search failed",
            "details": err.Error(),
        })
        return
    }

    response := models.SearchResponse{
        Documents: docs,
        Total:     total,
        Query:     req.Query,
        Took:      0, // Optionally set if you measure time
    }
    c.JSON(http.StatusOK, response)
}

// Health handles health check requests
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