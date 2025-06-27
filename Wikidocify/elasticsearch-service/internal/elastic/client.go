package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"wikidocify/elasticsearch-service/internal/models"

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

    // Mapping matches SearchDocument fields
    mapping := map[string]interface{}{
        "mappings": map[string]interface{}{
            "properties": map[string]interface{}{
                "id":        map[string]interface{}{"type": "integer"},
                "title":     map[string]interface{}{"type": "text", "analyzer": "standard"},
                "content":   map[string]interface{}{"type": "text", "analyzer": "standard"},
                "author":    map[string]interface{}{"type": "keyword"},
                "created_at": map[string]interface{}{"type": "date"},
                "updated_at": map[string]interface{}{"type": "date"},
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
    return nil
}

// IndexDocument indexes a SearchDocument in Elasticsearch
func (c *Client) IndexDocument(doc *models.SearchDocument) error {
    docJSON, err := json.Marshal(doc)
    if err != nil {
        return err
    }
    req := esapi.IndexRequest{
        Index:      c.index,
        DocumentID: fmt.Sprint(doc.ID),
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

// DeleteDocument deletes a document by ID from Elasticsearch
func (c *Client) DeleteDocument(id uint32) error {
    req := esapi.DeleteRequest{
        Index:      c.index,
        DocumentID: fmt.Sprint(id),
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

// Search performs a search query with filters and pagination
func (c *Client) Search(req *models.SearchRequest) ([]models.SearchDocument, int64, error) {
    queryFields := []string{}
    switch req.Type {
    case "title":
        queryFields = []string{"title"}
    case "content":
        queryFields = []string{"content"}
    default:
        queryFields = []string{"title^2", "content"}
    }

    esQuery := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []interface{}{
                    map[string]interface{}{
                        "multi_match": map[string]interface{}{
                            "query":  req.Query,
                            "fuzziness": "AUTO",
                            "fields": queryFields,
                        },
                    },
                },
            },
        },
        "from": req.Offset,
        "size": req.Limit,
    }

    // Optional author filter
    if req.Author != "" {
        boolQuery := esQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})
        boolQuery["filter"] = []interface{}{
            map[string]interface{}{
                "term": map[string]interface{}{
                    "author": req.Author,
                },
            },
        }
    }

    queryJSON, err := json.Marshal(esQuery)
    if err != nil {
        return nil, 0, err
    }

    start := time.Now()
    res, err := c.es.Search(
        c.es.Search.WithContext(context.Background()),
        c.es.Search.WithIndex(c.index),
        c.es.Search.WithBody(bytes.NewReader(queryJSON)),
    )
    if err != nil {
        return nil, 0, err
    }
    defer res.Body.Close()
    if res.IsError() {
        return nil, 0, fmt.Errorf("search failed: %s", res.Status())
    }

    var result map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
        return nil, 0, err
    }
    hits, ok := result["hits"].(map[string]interface{})
    if !ok {
        return nil, 0, fmt.Errorf("invalid search response format")
    }
    hitsList, _ := hits["hits"].([]interface{})
    docs := make([]models.SearchDocument, 0, len(hitsList))
    for _, hit := range hitsList {
        hitMap, _ := hit.(map[string]interface{})
        source, _ := hitMap["_source"].(map[string]interface{})
        var doc models.SearchDocument
        // Map fields from ES source to SearchDocument
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
        docs = append(docs, doc)
    }
    total := int64(0)
    if v, ok := hits["total"].(map[string]interface{}); ok {
        if val, ok := v["value"].(float64); ok {
            total = int64(val)
        }
    }
    _ = time.Since(start) // You can use this for took_ms if needed
    return docs, total, nil
}

func (c *Client) HealthCheck() error {
    return c.ping()
}