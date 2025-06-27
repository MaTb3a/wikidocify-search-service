package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// DocEvent matches the event structure expected by the consumer
type DocEvent struct {
	Event   string `json:"event"`   // "created", "updated", "deleted"
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	// Add more fields if needed (e.g., Author, CreatedAt, etc.)
}

var KafkaWriter *kafka.Writer

func InitKafkaWriter() {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	if broker == "" || topic == "" {
		log.Fatal("[KAFKA] KAFKA_BROKER and KAFKA_TOPIC must be set")
	}
	KafkaWriter = &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
		BatchTimeout: 10 * time.Millisecond,
	}
	log.Printf("[KAFKA] Producer initialized for broker %s, topic %s", broker, topic)
}

func PublishDocEvent(eventType, id, title, content string) error {
	if KafkaWriter == nil {
		log.Println("[KAFKA] KafkaWriter is not initialized")
		return nil
	}
	event := DocEvent{
		Event:   eventType,
		ID:      id,
		Title:   title,
		Content: content,
	}
	value, err := json.Marshal(event)
	if err != nil {
		log.Printf("[KAFKA] Failed to marshal event: %v", err)
		return err
	}
	msg := kafka.Message{
		Key:   []byte(id),
		Value: value,
		Time:  time.Now(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = KafkaWriter.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("[KAFKA] Failed to publish event: %v", err)
	}
	return err
}