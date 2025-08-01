package consumer

import (
	"interview-service/adapter/inbound/kafka/abstract"
)

type KafkaConsumer struct {
	handler abstract.HandleMessageInterface
}

// NewKafkaConsumer MOCK - Править не нужно
func NewKafkaConsumer(handler abstract.HandleMessageInterface) (*KafkaConsumer, error) {
	return &KafkaConsumer{handler: handler}, nil
}

func (c *KafkaConsumer) StartConsuming() error {
	return nil
}

func (c *KafkaConsumer) Close() error {
	return nil
}

func (c *KafkaConsumer) GetTopic(topic string) string {
	return topic
}

func (c *KafkaConsumer) GetTopics(topics []string) []string {
	return topics
}
