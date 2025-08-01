package abstract

import (
	"context"

	"github.com/IBM/sarama"
)

// HandleMessageInterface Обработка сообщение по одному
type HandleMessageInterface interface {

	// Handle Обрабатывает сообщения по одному
	Handle(ctx context.Context, message *sarama.ConsumerMessage) error

	// GetTopics Отдаёт прослушиваемые топики, без префикса
	GetTopics() []string

	// MarkMessageMode Режим для коммитов
	MarkMessageMode() string
}
