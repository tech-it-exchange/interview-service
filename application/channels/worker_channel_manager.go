package channels

import (
	"github.com/google/uuid"
	"interview-service/infrastructure/config"
	"log/slog"
	"sync"
)

type WorkerChannelManager struct {
	balanceChannel map[uuid.UUID]chan struct{}

	channelConfig config.ChannelConfig
	logger        *slog.Logger

	mutex sync.RWMutex
}

func NewWorkerChannelManager(
	channelConfig config.ChannelConfig,
	logger *slog.Logger,
) *WorkerChannelManager {
	return &WorkerChannelManager{
		balanceChannel: make(map[uuid.UUID]chan struct{}, channelConfig.BalanceChannelBufferSize),

		channelConfig: channelConfig,
		logger:        logger,
	}
}

func (m *WorkerChannelManager) GetBalanceChannel(instrumentId uuid.UUID) chan struct{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	channel, ok := m.balanceChannel[instrumentId]
	if !ok {
		channel = make(chan struct{}, m.channelConfig.BalanceChannelBufferSize)
		m.balanceChannel[instrumentId] = channel
	}

	return channel
}

func (m *WorkerChannelManager) TriggerBalance(instrumentId uuid.UUID) {
	m.mutex.RLock()
	channel, ok := m.balanceChannel[instrumentId]
	m.mutex.RUnlock()

	if !ok {
		m.logger.Error("Не найден триггер канал balanceChannel по instrumentId", instrumentId)

		return
	}

	select {
	case channel <- struct{}{}:
	default:
		// если уже есть сигнал — ничего не делаем
	}
}
