package service

import (
	"interview-service/adapter/outbound/localStorage"
	"log/slog"
)

type Services struct {
	Kafka *KafkaService
}

func NewServices(
	logger *slog.Logger,
	storages *localStorage.StorageContainer,
) *Services {
	kafkaService := NewKafkaService(
		storages.Instrument,
		storages.Kafka,
		logger,
	)

	return &Services{
		Kafka: kafkaService,
	}
}
