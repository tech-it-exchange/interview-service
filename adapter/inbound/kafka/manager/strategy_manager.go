package manager

import (
	"interview-service/adapter/inbound/kafka/manager/abstract"
	"log"
)

type strategyManager struct {
	handlers map[string]abstract.EventStrategyInterface
	logger   *log.Logger
}

func NewManager(
	logger *log.Logger,
) StrategyManagerInterface {

	handlers := map[string]abstract.EventStrategyInterface{}

	return &strategyManager{
		logger:   logger,
		handlers: handlers,
	}
}

// GetHandler Возвращает обработчик
func (m *strategyManager) GetHandler(contractName string) abstract.EventStrategyInterface {
	handler, ok := m.handlers[contractName]
	if !ok {
		m.logger.Printf("не найден обработчик для контракта: %s", contractName)
		return nil
	}

	return handler
}
