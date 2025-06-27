package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"wikidocify-search-service/internal/events"
	kafkaHandlers "wikidocify-search-service/internal/kafka-handlers"
	documenthandlers "wikidocify-search-service/internal/kafka-handlers/document-handlers"
)

type DocumentConsumer struct {
	Broker string
	Topic  string
}

func (c *DocumentConsumer) getHandlerMap() map[string]kafkaHandlers.Handler {
	return map[string]kafkaHandlers.Handler{
		"create": &documenthandlers.CreateHandler{},
		"update": &documenthandlers.UpdateHandler{},
		"delete": &documenthandlers.DeleteHandler{},
	}
}

func (c *DocumentConsumer) Consume(ctx context.Context) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{c.Broker},
		Topic:   c.Topic,
		GroupID: "document-group",
	})
	defer reader.Close()

	handlerMap := c.getHandlerMap()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		var event events.DocumentEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			fmt.Println("Invalid JSON:", err)
			continue
		}

		handler := handlerMap[event.EventType]
		if handler == nil {
			fmt.Println("No handler for type:", event.EventType)
			continue
		}

		if err := handler.Handle(ctx, event); err != nil {
			fmt.Println("Handler error:", err)
		}
	}
}
func init() {
	Register("document-topic", &DocumentConsumer{
		Broker: "localhost:9092",
		Topic:  "document-topic",
	})
}
