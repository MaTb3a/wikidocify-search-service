package kafkaHandlers

import (
	"context"
	"fmt"
)

type CreateHandler struct {
}

func (h *CreateHandler) Handle(ctx context.Context, data any) error {
	fmt.Println("AddHandler processing:", data)
	// your logic here
	return nil
}
