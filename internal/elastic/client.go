// internal/elastic/client.go
package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"wikidocify-search-service/internal/models"
)

type Client struct {
	es    *elasticsearch.Client
	index string
}

func NewClient(url, username, password, index string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	client := &Client{
		es:    es,
		index: index,
	}

	// Test connection
	if err := client.ping(); err != nil {
		return nil, fmt.Errorf("elasticsearch connection failed: %w", err)
	}

	// Create index if it doesn't exist
	if err := client.createIndexIfNotExists(); err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	return client, nil
}

func (c *Client) ping() error {
	res, err := c.es.Info()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch ping failed: %s", res.Status())
	}

	return nil
}

func (c *Client) createIndexIfNotExists() error {
	// Check if index exists
	res, err := c.es.Indices.Exists([]string{c.index})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		// Index exists
		return nil
	}

	// Create index with mapping
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "integer",
				},
				"title": map[string]interface{}{
					"type": "text",
					"analyzer": "standard",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"content": map[string]interface{}{
					"type": "text",
					"analyzer": "standard",
				},
				"author": map[string]interface{}{
					"type": "keyword",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
				"updated_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}

	mappingJSON, err := json.Marshal(mapping)
	if err != nil {
		return err
	}

	req := esapi.IndicesCreateRequest{
		Index: c.index,
		Body:  bytes.NewReader(mappingJSON),
	}

	res, err = req.Do(context.Background(), c.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index: %s", res.Status())
	}

	log.Printf("Created Elasticsearch index: %s", c.index)
	return nil
}

func (c *Client) IndexDocument(doc *models.SearchDocument) error {
	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      c.index,
		DocumentID: fmt.Sprintf("%d", doc.ID),
		Body:       bytes.NewReader(docJSON),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to index document: %s", res.Status())
	}

	return nil
}

func (c *Client) BulkIndex(docs []*models.SearchDocument) error {
	if len(docs) == 0 {
		return nil
	}

	var buf bytes.Buffer

	for _, doc := range docs {
		// Index action
		action := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": c.index,
				"_id":    fmt.Sprintf("%d", doc.ID),
			},
		}
		actionJSON, _ := json.Marshal(action)
		buf.Write(actionJSON)
		buf.WriteByte('\n')

		// Document
		docJSON, _ := json.Marshal(doc)
		buf.Write(docJSON)
		buf.WriteByte('\n')
	}

	req := esapi.BulkRequest{
		Body:    &buf,
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk index failed: %s", res.Status())
	}

	return nil
}

func (c *Client) Search(searchReq *models.SearchRequest) (*models.SearchResponse, error) {
	query := c.buildQuery(searchReq)
	
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{c.index},
		Body:  bytes.NewReader(queryJSON),
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search failed: %s", res.Status())
	}

	var searchResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, err
	}

	return c.parseSearchResponse(searchResult, searchReq.Query)
}

func (c *Client) buildQuery(searchReq *models.SearchRequest) map[string]interface{} {
	// Set defaults
	if searchReq.Limit <= 0 || searchReq.Limit > 100 {
		searchReq.Limit = 10
	}
	if searchReq.Offset < 0 {
		searchReq.Offset = 0
	}
	if searchReq.Type == "" {
		searchReq.Type = "all"
	}

	var should []map[string]interface{}

	switch searchReq.Type {
	case "title":
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"title": map[string]interface{}{
					"query": searchReq.Query,
					"boost": 2.0,
				},
			},
		})
	case "content":
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"content": searchReq.Query,
			},
		})
	default: // "all"
		should = append(should,
			map[string]interface{}{
				"match": map[string]interface{}{
					"title": map[string]interface{}{
						"query": searchReq.Query,
						"boost": 2.0,
					},
				},
			},
			map[string]interface{}{
				"match": map[string]interface{}{
					"content": searchReq.Query,
				},
			},
		)
	}

	boolQuery := map[string]interface{}{
		"should": should,
		"minimum_should_match": 1,
	}

	// Add author filter if specified
	if searchReq.Author != "" {
		boolQuery["filter"] = []map[string]interface{}{
			{
				"term": map[string]interface{}{
					"author": searchReq.Author,
				},
			},
		}
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": boolQuery,
		},
		"from": searchReq.Offset,
		"size": searchReq.Limit,
		"sort": []map[string]interface{}{
			{"_score": map[string]string{"order": "desc"}},
			{"updated_at": map[string]string{"order": "desc"}},
		},
	}

	return query
}

func (c *Client) parseSearchResponse(result map[string]interface{}, query string) (*models.SearchResponse, error) {
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid search response format")
	}

	total, _ := hits["total"].(map[string]interface{})
	totalValue, _ := total["value"].(float64)

	took, _ := result["took"].(float64)

	hitsList, _ := hits["hits"].([]interface{})
	documents := make([]models.SearchDocument, 0, len(hitsList))

	for _, hit := range hitsList {
		hitMap, _ := hit.(map[string]interface{})
		source, _ := hitMap["_source"].(map[string]interface{})

		doc := models.SearchDocument{}
		if id, ok := source["id"].(float64); ok {
			doc.ID = uint32(id)
		}
		if title, ok := source["title"].(string); ok {
			doc.Title = title
		}
		if content, ok := source["content"].(string); ok {
			doc.Content = content
		}
		if author, ok := source["author"].(string); ok {
			doc.Author = author
		}
		if createdAt, ok := source["created_at"].(string); ok {
			if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
				doc.CreatedAt = t
			}
		}
		if updatedAt, ok := source["updated_at"].(string); ok {
			if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
				doc.UpdatedAt = t
			}
		}

		documents = append(documents, doc)
	}

	return &models.SearchResponse{
		Documents: documents,
		Total:     int64(totalValue),
		Took:      int64(took),
		Query:     query,
	}, nil
}

func (c *Client) DeleteDocument(docID uint32) error {
	req := esapi.DeleteRequest{
		Index:      c.index,
		DocumentID: fmt.Sprintf("%d", docID),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("failed to delete document: %s", res.Status())
	}

	return nil
}

func (c *Client) HealthCheck() error {
	return c.ping()
}