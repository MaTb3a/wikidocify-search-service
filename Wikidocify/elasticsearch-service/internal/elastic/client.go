package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Client struct {
	es    *elasticsearch.Client
	index string
}

func NewClient(url, index string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
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
	res, err := c.es.Indices.Exists([]string{c.index})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		return nil // Index exists
	}

	// Minimal mapping for search service
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{"type": "keyword"},
				"title": map[string]interface{}{
					"type":     "text",
					"analyzer": "standard",
				},
				"content": map[string]interface{}{
					"type":     "text",
					"analyzer": "standard",
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

func (c *Client) IndexDocument(id, title, content string) error {
	doc := map[string]interface{}{
		"id":      id,
		"title":   title,
		"content": content,
	}
	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	req := esapi.IndexRequest{
		Index:      c.index,
		DocumentID: id,
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

func (c *Client) DeleteDocument(id string) error {
	req := esapi.DeleteRequest{
		Index:      c.index,
		DocumentID: id,
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

func (c *Client) Search(query string, limit, offset int) ([]map[string]interface{}, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title^2", "content"},
			},
		},
		"from": offset,
		"size": limit,
	}
	queryJSON, err := json.Marshal(esQuery)
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
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid search response format")
	}
	hitsList, _ := hits["hits"].([]interface{})
	docs := make([]map[string]interface{}, 0, len(hitsList))
	for _, hit := range hitsList {
		hitMap, _ := hit.(map[string]interface{})
		source, _ := hitMap["_source"].(map[string]interface{})
		docs = append(docs, source)
	}
	return docs, nil
}

func (c *Client) HealthCheck() error {
	return c.ping()
}
