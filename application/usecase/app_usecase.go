package usecase

import (
	"context"
	"github.com/google/uuid"
	"interview-service/adapter/inbound/fasthttp/dto"
	"interview-service/adapter/outbound/postgres/repositories"
	"interview-service/domain"
	"interview-service/port/outbound/localStorage"
	"log/slog"
)

type AppUseCase struct {
	instrumentStorage localStorage.InstrumentStorageInterface
	balanceStorage    localStorage.BalanceStorageInterface
	kafkaStorage      localStorage.KafkaStorageInterface
	userRepository    *repositories.UserRepository
	logger            *slog.Logger
}

func NewAppUseCase(
	instrumentStorage localStorage.InstrumentStorageInterface,
	balanceStorage localStorage.BalanceStorageInterface,
	userRepository *repositories.UserRepository,
	logger *slog.Logger,
) *AppUseCase {
	return &AppUseCase{
		instrumentStorage: instrumentStorage,
		balanceStorage:    balanceStorage,
		userRepository:    userRepository,
		logger:            logger,
	}
}

func (u *AppUseCase) GetLoss(ctx context.Context) (int64, error) {
	return -1000, nil
}

func (u *AppUseCase) CreateOrder(ctx context.Context, createOrderRequestDto *dto.CreateOrderRequestDto) (uuid.UUID, error) {
	// Создание ордера ...
	return uuid.New(), nil
}

func (u *AppUseCase) CreateInstrument(ctx context.Context, instrument *domain.InstrumentEntity) (uuid.UUID, error) {
	u.instrumentStorage.SaveInstrument(instrument)
	u.balanceStorage.SaveInstrumentBalance(instrument.InstrumentId)

	return uuid.New(), nil
}

func (u *AppUseCase) GetUserName(ctx context.Context, name string) string {
	return u.userRepository.GetUserByName(name)
}

func (u *AppUseCase) GetInstrumentBalance(ctx context.Context, instrumentId uuid.UUID) string {
	return u.balanceStorage.GetInstrumentBalance(instrumentId).String()
}

func (u *AppUseCase) GetActiveConsumerTopics() []string {
	consumers := u.kafkaStorage.GetKafkaConsumers()
	var topics []string

	for topic := range consumers {
		topics = append(topics, topic)
	}

	return topics
}
