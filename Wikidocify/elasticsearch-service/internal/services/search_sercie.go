// internal/services/search_service.go
package services

import (
	"fmt"
	"log"
	"time"

	"wikidocify/elasticsearch-service/internal/elastic"
	"wikidocify/elasticsearch-service/internal/models"
)

type SearchService struct {
	esClient     *elastic.Client
	docService   *DocServiceClient
	syncInterval time.Duration
	batchSize    int
	enableSync   bool
	lastSyncTime time.Time
}

func NewSearchService(esClient *elastic.Client, docService *DocServiceClient, syncInterval time.Duration, batchSize int, enableSync bool) *SearchService {
	return &SearchService{
		esClient:     esClient,
		docService:   docService,
		syncInterval: syncInterval,
		batchSize:    batchSize,
		enableSync:   enableSync,
	}
}

// Search performs a search using the ES client and returns docs and total count
func (s *SearchService) Search(req *models.SearchRequest) ([]models.SearchDocument, int64, error) {
	return s.esClient.Search(req)
}

// SyncDocument fetches a document from doc service and indexes it in ES
func (s *SearchService) SyncDocument(docID uint32) error {
	// Get document from doc service
	doc, err := s.docService.GetDocument(docID)
	if err != nil {
		return fmt.Errorf("failed to get document from doc service: %w", err)
	}

	// Convert to search document
	searchDoc := doc.ToSearchDocument()

	// Index in Elasticsearch
	if err := s.esClient.IndexDocument(searchDoc); err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}

	log.Printf("Successfully synced document ID: %d", docID)
	return nil
}

// DeleteDocument removes a document from ES by ID
func (s *SearchService) DeleteDocument(docID uint32) error {
	if err := s.esClient.DeleteDocument(docID); err != nil {
		return fmt.Errorf("failed to delete document from Elasticsearch: %w", err)
	}

	log.Printf("Successfully deleted document ID: %d", docID)
	return nil
}

// FullSync fetches all documents from doc service and bulk indexes them in ES
func (s *SearchService) FullSync() error {
	log.Println("Starting full sync...")
	
	page := 1
	totalSynced := 0

	for {
		// Get documents in batches
		docs, err := s.docService.GetDocumentsPaginated(page, s.batchSize)
		if err != nil {
			return fmt.Errorf("failed to get documents from doc service: %w", err)
		}

		if len(docs) == 0 {
			break
		}

		// Convert to search documents
		searchDocs := make([]*models.SearchDocument, len(docs))
		for i, doc := range docs {
			searchDocs[i] = doc.ToSearchDocument()
		}

		// If you have BulkIndex implemented, use it. Otherwise, index one by one:
		for _, sd := range searchDocs {
			if err := s.esClient.IndexDocument(sd); err != nil {
				return fmt.Errorf("failed to index document: %w", err)
			}
		}

		totalSynced += len(docs)
		log.Printf("Synced batch %d: %d documents (total: %d)", page, len(docs), totalSynced)

		// If we got fewer documents than the batch size, we're done
		if len(docs) < s.batchSize {
			break
		}

		page++
	}

	s.lastSyncTime = time.Now()
	log.Printf("Full sync completed. Total documents synced: %d", totalSynced)
	return nil
}

// StartPeriodicSync runs FullSync on a schedule if enabled
func (s *SearchService) StartPeriodicSync() {
	if !s.enableSync {
		log.Println("Periodic sync is disabled")
		return
	}

	log.Printf("Starting periodic sync every %v", s.syncInterval)
	
	ticker := time.NewTicker(s.syncInterval)
	go func() {
		for range ticker.C {
			log.Println("Starting periodic sync...")
			if err := s.FullSync(); err != nil {
				log.Printf("Periodic sync failed: %v", err)
			} else {
				log.Println("Periodic sync completed successfully")
			}
		}
	}()
}

// GetSyncStatus returns sync status info
func (s *SearchService) GetSyncStatus() map[string]interface{} {
	return map[string]interface{}{
		"sync_enabled":    s.enableSync,
		"sync_interval":   s.syncInterval.String(),
		"batch_size":      s.batchSize,
		"last_sync_time":  s.lastSyncTime,
	}
}

// HealthCheck checks the health of ES and doc service
func (s *SearchService) HealthCheck() map[string]bool {
	status := map[string]bool{
		"elasticsearch": true,
		"doc_service":   true,
	}

	// Check Elasticsearch
	if err := s.esClient.HealthCheck(); err != nil {
		status["elasticsearch"] = false
		log.Printf("Elasticsearch health check failed: %v", err)
	}

	// Check Doc Service
	if err := s.docService.HealthCheck(); err != nil {
		status["doc_service"] = false
		log.Printf("Doc service health check failed: %v", err)
	}

	return status
}