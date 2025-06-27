package consumer

import (
	"context"
)

type ConsumerStrategy interface {
	Consume(ctx context.Context) error
}
