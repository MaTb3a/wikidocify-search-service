// internal/services/doc_service.go
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"wikidocify-search-service/internal/models"
)

type DocServiceClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewDocServiceClient(baseURL, apiKey string, timeout time.Duration) *DocServiceClient {
	return &DocServiceClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *DocServiceClient) GetAllDocuments() ([]*models.Document, error) {
	url := fmt.Sprintf("%s/api/v1/documents", c.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("doc service returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Documents []*models.Document `json:"documents"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Documents, nil
}

func (c *DocServiceClient) GetDocumentsPaginated(page, limit int) ([]*models.Document, error) {
	url := fmt.Sprintf("%s/api/v1/documents?page=%d&limit=%d", c.baseURL, page, limit)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("doc service returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Documents []*models.Document `json:"documents"`
		Total     int               `json:"total"`
		Page      int               `json:"page"`
		Limit     int               `json:"limit"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Documents, nil
}

func (c *DocServiceClient) GetDocument(id uint32) (*models.Document, error) {
	url := fmt.Sprintf("%s/api/v1/documents/%d", c.baseURL, id)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("document not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("doc service returned status %d: %s", resp.StatusCode, string(body))
	}

	var document models.Document
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &document, nil
}

func (c *DocServiceClient) HealthCheck() error {
	url := fmt.Sprintf("%s/health", c.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("doc service health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("doc service health check failed with status: %d", resp.StatusCode)
	}

	return nil
}