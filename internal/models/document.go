// internal/models/document.go
package models

import (
	"time"
)

// Document represents the original document structure from the doc service
type Document struct {
	ID        uint32    `json:"id"`
	Title     string    `json:"title"`
	Content   []byte    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SearchDocument represents the document structure in Elasticsearch
type SearchDocument struct {
	ID        uint32    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToSearchDocument converts Document to SearchDocument
func (d *Document) ToSearchDocument() *SearchDocument {
	return &SearchDocument{
		ID:        d.ID,
		Title:     d.Title,
		Content:   string(d.Content),
		Author:    d.Author,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

// SearchRequest represents search request parameters
type SearchRequest struct {
	Query  string `json:"query" form:"query" binding:"required"`
	Type   string `json:"type" form:"type"`           // "title", "content", or "all"
	Limit  int    `json:"limit" form:"limit"`         // default 10
	Offset int    `json:"offset" form:"offset"`       // default 0
	Author string `json:"author" form:"author"`       // optional filter
}

// SearchResponse represents search response
type SearchResponse struct {
	Documents []SearchDocument `json:"documents"`
	Total     int64            `json:"total"`
	Took      int64            `json:"took_ms"`
	Query     string           `json:"query"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status          string            `json:"status"`
	Timestamp       time.Time         `json:"timestamp"`
	ElasticsearchOK bool              `json:"elasticsearch_ok"`
	DocServiceOK    bool              `json:"doc_service_ok"`
	Details         map[string]string `json:"details,omitempty"`
}