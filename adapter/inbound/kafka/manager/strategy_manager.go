package manager

import (
	"interview-service/adapter/inbound/kafka/manager/abstract"
	"log/slog"
)

type strategyManager struct {
	handlers map[string]abstract.EventStrategyInterface
	logger   *slog.Logger
}

func NewManager(
	logger *slog.Logger,
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
		m.logger.Warn("не найден обработчик для контракта", "contractName", contractName)
		return nil
	}

	return handler
}
