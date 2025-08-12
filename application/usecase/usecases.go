package usecase

import (
	"interview-service/adapter/outbound/localStorage"
	"interview-service/adapter/outbound/postgres"
	"log/slog"
)

type UseCases struct {
	Worker *AppUseCase
}

func NewUseCases(
	storageContainer *localStorage.StorageContainer,
	repositories postgres.Repositories,
	logger *slog.Logger,
) *UseCases {
	calculateOrderUseCase := NewAppUseCase(
		storageContainer.Instrument,
		storageContainer.Balance,
		repositories.User,
		logger,
	)

	return &UseCases{
		Worker: calculateOrderUseCase,
	}
}
