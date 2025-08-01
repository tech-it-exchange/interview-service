package abstract

import (
	"context"
	"github.com/IBM/sarama"
)

type EventStrategyInterface interface {
	ProcessMessage(ctx context.Context, message *sarama.ConsumerMessage) error
}
