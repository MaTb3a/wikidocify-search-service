package kafkaHandlers

import (
	"context"
	"fmt"
)

type DeleteHandler struct{}

func (h *DeleteHandler) Handle(ctx context.Context, data any) error {
	fmt.Println("DeleteHandler processing:", data)
	// your logic here
	return nil
}
