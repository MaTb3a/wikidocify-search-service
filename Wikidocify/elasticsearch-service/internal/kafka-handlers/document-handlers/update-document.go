package kafkaHandlers

import (
	"context"
	"fmt"
)

type UpdateHandler struct{}

func (h *UpdateHandler) Handle(ctx context.Context, data any) error {
	fmt.Println("UpdateHandler processing:", data)
	return nil
}
