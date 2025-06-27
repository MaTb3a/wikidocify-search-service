// shared/models/kafka_message.go
package models

import (
    "encoding/json"
    "time"
)

// EventType represents different types of document events
type EventType string

const (
    DocumentCreated EventType = "document.created"
    DocumentUpdated EventType = "document.updated"
    DocumentDeleted EventType = "document.deleted"
)

// DocumentEvent represents a document event message
type DocumentEvent struct {
    ID          string                 `json:"id"`
    EventType   EventType              `json:"event_type"`
    Timestamp   time.Time              `json:"timestamp"`
    DocumentID  string                 `json:"document_id"`
    Data        map[string]interface{} `json:"data,omitempty"`
    Metadata    map[string]string      `json:"metadata,omitempty"`
}

// ToJSON converts DocumentEvent to JSON bytes
func (de *DocumentEvent) ToJSON() ([]byte, error) {
    return json.Marshal(de)
}

// FromJSON creates DocumentEvent from JSON bytes
func (de *DocumentEvent) FromJSON(data []byte) error {
    return json.Unmarshal(data, de)
}

// NewDocumentEvent creates a new document event
func NewDocumentEvent(eventType EventType, documentID string, data map[string]interface{}) *DocumentEvent {
    return &DocumentEvent{
        ID:         generateEventID(),
        EventType:  eventType,
        Timestamp:  time.Now(),
        DocumentID: documentID,
        Data:       data,
        Metadata:   make(map[string]string),
    }
}

// generateEventID generates a unique event ID
func generateEventID() string {
    // In a real implementation, use UUID or similar
    return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}