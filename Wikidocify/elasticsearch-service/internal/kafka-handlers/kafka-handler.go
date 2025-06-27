package kafkaHandlers

import (
	"context"
)

// Strategy Interface
type Handler interface {
	Handle(ctx context.Context, data any) error
}
