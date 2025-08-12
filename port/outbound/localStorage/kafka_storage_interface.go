package localStorage

import (
	"github.com/google/uuid"
	"interview-service/adapter/inbound/kafka/consumer"
)

type KafkaStorageInterface interface {
	SaveKafkaConsumer(topic string, consumer *consumer.KafkaConsumer)
	SaveKafkaConsumersMap(consumersMap map[string]*consumer.KafkaConsumer)
	SaveActiveSpotInstrumentsMap(activeSpotInstrumentsMap map[uuid.UUID]bool)
	GetKafkaConsumers() map[string]*consumer.KafkaConsumer
	HasActiveSpotInstrument(instrumentId uuid.UUID) bool
	HasKafkaConsumer(topic string) bool
	SaveActiveSpotInstrument(instrumentId uuid.UUID)
	SetSpotNeedInstrumentDispatch(value bool)
	IsNeedSpotInstrumentDispatch() bool
}
