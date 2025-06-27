// internal/kafka/consumer.go
package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"wikidocify/elasticsearch-service/internal/services"

	"github.com/segmentio/kafka-go"
)

type DocEvent struct {
	Event   string `json:"event"`   // "created", "updated", "deleted"
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func StartConsumer(searchService *services.SearchService) {
	broker := os.Getenv("KAFKA_BROKER") // e.g., "kafka:9092"
	topic := os.Getenv("KAFKA_TOPIC")   // e.g., "document-events"
	if broker == "" || topic == "" {
		log.Fatal("KAFKA_BROKER and KAFKA_TOPIC must be set")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{broker},
		Topic:       topic,
		GroupID:     "search-service-group",
		MinBytes:    1,    // 1B
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
	})

	log.Println("🔄 Kafka consumer started")
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		var ev DocEvent
		if err := json.Unmarshal(msg.Value, &ev); err != nil {
			log.Printf("Failed to unmarshal Kafka event: %v", err)
			continue
		}
		log.Printf("Received Kafka event: %+v", ev)

		// Parse ID as uint32
		id, err := strconv.ParseUint(ev.ID, 10, 32)
		if err != nil {
			log.Printf("Invalid document ID in event: %v", err)
			continue
		}

		switch ev.Event {
		case "created", "updated":
			if err := searchService.SyncDocument(uint32(id)); err != nil {
				log.Printf("Failed to sync document: %v", err)
			}
		case "deleted":
			if err := searchService.DeleteDocument(uint32(id)); err != nil {
				log.Printf("Failed to delete document: %v", err)
			}
		default:
			log.Printf("Unknown event type: %s", ev.Event)
		}
	}
}
