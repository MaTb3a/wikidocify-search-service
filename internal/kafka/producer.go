package kafka

import (
	"context"
	"log"
	"os"
	"github.com/segmentio/kafka-go"
)

var (
	KafkaWriter *kafka.Writer
)

func InitKafkaWriter() {
	broker := os.Getenv("KAFKA_BROKER") // e.g., "kafka:9092"
	topic := os.Getenv("KAFKA_TOPIC")   // e.g., "document-events"
	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("✅ Kafka writer initialized")
}

func PublishDocEvent(eventType string, docID string, title string, content string) error {
	msg := kafka.Message{
		Key: []byte(docID),
		Value: []byte(`{
			"event":"` + eventType + `",
			"id":"` + docID + `",
			"title":"` + title + `",
			"content":"` + content + `"
		}`),
	}
	return KafkaWriter.WriteMessages(context.Background(), msg)
}
