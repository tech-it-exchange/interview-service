package localStorage

import (
	"github.com/google/uuid"
	"interview-service/adapter/inbound/kafka/consumer"
	"interview-service/port/outbound/localStorage"
	"sync"
)

type kafkaStorage struct {
	kafkaConsumersMap          map[string]*consumer.KafkaConsumer // Мапа потребителей кафка, ключом является собранный topic
	activeSpotInstrumentsMap   map[uuid.UUID]bool                 // Мапа активных инструментов с которых были созданы потребители
	needSpotInstrumentDispatch bool
	mu                         sync.RWMutex
}

func NewKafkaStorage() localStorage.KafkaStorageInterface {
	return &kafkaStorage{
		kafkaConsumersMap:        make(map[string]*consumer.KafkaConsumer),
		activeSpotInstrumentsMap: make(map[uuid.UUID]bool),
	}
}

func (s *kafkaStorage) SaveKafkaConsumer(topic string, consumer *consumer.KafkaConsumer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.kafkaConsumersMap[topic] = consumer
}

func (s *kafkaStorage) SaveKafkaConsumersMap(consumersMap map[string]*consumer.KafkaConsumer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.kafkaConsumersMap = consumersMap
}

func (s *kafkaStorage) SaveActiveSpotInstrumentsMap(activeSpotInstrumentsMap map[uuid.UUID]bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.activeSpotInstrumentsMap = activeSpotInstrumentsMap
}

func (s *kafkaStorage) GetKafkaConsumers() map[string]*consumer.KafkaConsumer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	consumersMapCopy := make(map[string]*consumer.KafkaConsumer, len(s.kafkaConsumersMap))
	for key, kafkaConsumer := range s.kafkaConsumersMap {
		consumersMapCopy[key] = kafkaConsumer
	}

	return consumersMapCopy
}

func (s *kafkaStorage) HasActiveSpotInstrument(instrumentId uuid.UUID) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.activeSpotInstrumentsMap[instrumentId]

	return ok
}

func (s *kafkaStorage) HasKafkaConsumer(topic string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.kafkaConsumersMap[topic]

	return ok
}

func (s *kafkaStorage) SaveActiveSpotInstrument(instrumentId uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.activeSpotInstrumentsMap[instrumentId] = true
}

func (s *kafkaStorage) SetSpotNeedInstrumentDispatch(value bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.needSpotInstrumentDispatch = value
}

func (s *kafkaStorage) IsNeedSpotInstrumentDispatch() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.needSpotInstrumentDispatch
}
