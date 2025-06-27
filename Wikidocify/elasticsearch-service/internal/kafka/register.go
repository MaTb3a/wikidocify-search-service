package consumer

import (
	"context"
)

// Create a shared map to hold all consumers
var registry = make(map[string]ConsumerStrategy)

// Function that other consumers will call to register themselves
func Register(topic string, strategy ConsumerStrategy) {
	registry[topic] = strategy
}

// Function to start all registered consumers
func StartAll(ctx context.Context) {
	for _, consumer := range registry {
		go consumer.Consume(ctx) // start each in background
	}
}
