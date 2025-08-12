package service

import (
	"fmt"
	"github.com/google/uuid"
	"interview-service/adapter/inbound/kafka/consumer"
	"interview-service/port/outbound/localStorage"
	"log/slog"
)

type KafkaService struct {
	instrumentStorage localStorage.InstrumentStorageInterface
	kafkaStorage      localStorage.KafkaStorageInterface
	logger            *slog.Logger
}

func NewKafkaService(
	instrumentStorage localStorage.InstrumentStorageInterface,
	kafkaStorage localStorage.KafkaStorageInterface,
	logger *slog.Logger,
) *KafkaService {
	return &KafkaService{
		instrumentStorage: instrumentStorage,
		kafkaStorage:      kafkaStorage,
		logger:            logger,
	}
}

func (s *KafkaService) SaveKafkaConsumer(topic string, consumer *consumer.KafkaConsumer) {
	s.kafkaStorage.SaveKafkaConsumer(topic, consumer)
}

func (s *KafkaService) SaveKafkaConsumersMap(consumersMap map[string]*consumer.KafkaConsumer) {
	s.kafkaStorage.SaveKafkaConsumersMap(consumersMap)
}

func (s *KafkaService) SaveActiveSpotInstrumentsMap(activeSpotInstrumentsMap map[uuid.UUID]bool) {
	s.kafkaStorage.SaveActiveSpotInstrumentsMap(activeSpotInstrumentsMap)
}

func (s *KafkaService) GetKafkaConsumers() map[string]*consumer.KafkaConsumer {
	return s.kafkaStorage.GetKafkaConsumers()
}

func (s *KafkaService) HasActiveSpotInstrument(instrumentId uuid.UUID) bool {
	return s.kafkaStorage.HasActiveSpotInstrument(instrumentId)
}

func (s *KafkaService) HasKafkaConsumer(topic string) bool {
	return s.kafkaStorage.HasKafkaConsumer(topic)
}

func (s *KafkaService) SaveActiveSpotInstrument(instrumentId uuid.UUID) {
	s.kafkaStorage.SaveActiveSpotInstrument(instrumentId)
}

func (s *KafkaService) SetSpotNeedInstrumentDispatch(value bool) {
	s.kafkaStorage.SetSpotNeedInstrumentDispatch(value)
}

func (s *KafkaService) IsNeedSpotInstrumentDispatch() bool {
	return s.kafkaStorage.IsNeedSpotInstrumentDispatch()
}

func (s *KafkaService) GetConfirmTopics() []string {
	return []string{"confirm.topic"}
}

func (s *KafkaService) GetCommandTopics() []string {
	return []string{
		"command.topic",
	}
}

func (s *KafkaService) GetSpotInstrumentTopics() map[uuid.UUID]string {
	topics := make(map[uuid.UUID]string)

	for _, instrumentId := range s.instrumentStorage.GetActiveSpotInstrumentIds() {
		topic := s.CreateSpotInstrumentTopic(instrumentId)
		topics[instrumentId] = topic
	}

	return topics
}

func (s *KafkaService) CreateSpotInstrumentTopic(instrumentId uuid.UUID) string {
	return fmt.Sprintf("instrument.topic.id.%s", instrumentId.String())
}
