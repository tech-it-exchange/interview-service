package usecase

import (
	"interview-service/adapter/outbound/localStorage"
	"interview-service/adapter/outbound/postgres"
	"log"
)

type UseCases struct {
	Worker *AppUseCase
}

func NewUseCases(
	storageContainer *localStorage.StorageContainer,
	repositories postgres.Repositories,
	logger *log.Logger,
) *UseCases {
	calculateOrderUseCase := NewAppUseCase(
		storageContainer.Instrument,
		storageContainer.Balance,
		storageContainer.Kafka,
		repositories.User,
		logger,
	)

	return &UseCases{
		Worker: calculateOrderUseCase,
	}
}
