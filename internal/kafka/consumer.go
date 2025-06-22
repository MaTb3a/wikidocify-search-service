// internal/kafka/consumer.go
package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/hossamhakim/wikidocify-search-service/internal/elastic"
)

type DocEvent struct {
	Event   string `json:"event"`
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func StartConsumer() {
	broker := os.Getenv("KAFKA_BROKER") // "kafka:9092"
	topic := os.Getenv("KAFKA_TOPIC")   // "document-events"
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     topic,
		GroupID:   "search-service-group",
		MinBytes:  1,    // 1B
		MaxBytes:  10e6, // 10MB
		StartOffset: kafka.LastOffset,
	})

	log.Println("🔄 Kafka consumer started")
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		var ev DocEvent
		if err := json.Unmarshal(msg.Value, &ev); err != nil {
			log.Printf("Failed to unmarshal Kafka event: %v", err)
			continue
		}
		log.Printf("Received Kafka event: %+v", ev)
		switch ev.Event {
		case "created", "updated":
			elastic.IndexDocument(ev.ID, ev.Title, ev.Content)
		case "deleted":
			elastic.DeleteDocument(ev.ID)
		}
	}
}
