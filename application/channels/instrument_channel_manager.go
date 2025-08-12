package channels

import (
	"github.com/google/uuid"
)

type InstrumentChannelManager struct {
	kafkaSpotInstrumentUuidChannel     chan uuid.UUID
	spotSyncEventInstrumentUuidChannel chan uuid.UUID
}

func NewInstrumentChannelManager() *InstrumentChannelManager {
	return &InstrumentChannelManager{
		kafkaSpotInstrumentUuidChannel:     make(chan uuid.UUID, 10),
		spotSyncEventInstrumentUuidChannel: make(chan uuid.UUID, 10),
	}
}

func (m *InstrumentChannelManager) GetKafkaSpotInstrumentIdChannel() chan uuid.UUID {
	return m.kafkaSpotInstrumentUuidChannel
}

func (m *InstrumentChannelManager) SendKafkaSpotInstrumentId(instrumentId uuid.UUID) {
	m.kafkaSpotInstrumentUuidChannel <- instrumentId
}

func (m *InstrumentChannelManager) GetSpotSyncEventInstrumentIdChannel() chan uuid.UUID {
	return m.spotSyncEventInstrumentUuidChannel
}

func (m *InstrumentChannelManager) SendSpotSyncEventSpotInstrumentId(instrumentId uuid.UUID) {
	m.spotSyncEventInstrumentUuidChannel <- instrumentId
}
