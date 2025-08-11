package handlers

import (
	"context"
	"interview-service/adapter/inbound/kafka/manager"
	"log"

	"github.com/IBM/sarama"
	libJson "github.com/tidwall/gjson"
)

type CommonHandler struct {
	strategyManager manager.StrategyManagerInterface
	topic           string
	messageMode     string
	logger          *log.Logger
}

func NewCommonHandler(
	strategyManager manager.StrategyManagerInterface,
	topic string,
	messageMode string,
	logger *log.Logger,
) *CommonHandler {
	return &CommonHandler{
		strategyManager: strategyManager,
		topic:           topic,
		messageMode:     messageMode,
		logger:          logger,
	}
}

// GetTopics Отдаёт прослушиваемые топики, без префикса
func (c *CommonHandler) GetTopics() []string {
	return []string{c.topic}
}

// Handle Обрабатывает сообщения по одному
func (c *CommonHandler) Handle(ctx context.Context, message *sarama.ConsumerMessage) error {
	handler := c.strategyManager.GetHandler(libJson.GetBytes(message.Value, "ContractName").String())
	if handler == nil {
		return nil
	}

	return handler.ProcessMessage(ctx, message)
}

// MarkMessageMode Режим для коммитов
func (c *CommonHandler) MarkMessageMode() string {
	return c.messageMode
}
