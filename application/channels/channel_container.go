package channels

import (
	"interview-service/infrastructure/config"
	"log"
)

type ChannelContainer struct {
	WorkerChannelManager     *WorkerChannelManager
	InstrumentChannelManager *InstrumentChannelManager
}

func NewChannelContainer(
	channelConfig config.ChannelConfig,
	logger *log.Logger,
) *ChannelContainer {
	workerChannelManager := NewWorkerChannelManager(channelConfig, logger)
	instrumentChannelManager := NewInstrumentChannelManager()

	return &ChannelContainer{
		WorkerChannelManager:     workerChannelManager,
		InstrumentChannelManager: instrumentChannelManager,
	}
}
