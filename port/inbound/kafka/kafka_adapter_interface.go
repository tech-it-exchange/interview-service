package kafka

import "context"

type KafkaAdapterInterface interface {
	InitConsumers() error
	StartConsuming() error
	CloseConsuming()
	ListenForNewInstrument(ctx context.Context)
}
