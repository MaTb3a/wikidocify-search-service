package kafka

import (
    "context"
    "log"
    "os"
    "github.com/segmentio/kafka-go"
    "encoding/json"

)

var KafkaWriter *kafka.Writer

func InitKafkaWriter() {
    broker := os.Getenv("KAFKA_BROKER")
    topic := os.Getenv("KAFKA_TOPIC")
    KafkaWriter = &kafka.Writer{
        Addr:     kafka.TCP(broker),
        Topic:    topic,
        Balancer: &kafka.LeastBytes{},
    }
    log.Println("[KAFKA] Producer initialized")
}


func PublishDocEvent(eventType, id, title, content string) error {
    event := map[string]string{
        "event":   eventType,
        "id":      id,
        "title":   title,
        "content": content,
    }
    value, err := json.Marshal(event)
    if err != nil {
        return err
    }
    msg := kafka.Message{
        Key:   []byte(id),
        Value: value,
    }
    return KafkaWriter.WriteMessages(context.Background(), msg)
}