package usecase

import (
	"context"
	"fmt"
	"interview-service/adapter/inbound/fasthttp/dto"
	"interview-service/adapter/outbound/postgres/repositories"
	"interview-service/domain"
	"interview-service/port/outbound/localStorage"
	"log"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AppUseCase struct {
	instrumentStorage localStorage.InstrumentStorageInterface
	balanceStorage    localStorage.BalanceStorageInterface
	kafkaStorage      localStorage.KafkaStorageInterface
	userRepository    *repositories.UserRepository
	// используем стандартный логгер для совместимости со старыми версиями Go
	logger *log.Logger
}

func NewAppUseCase(
	instrumentStorage localStorage.InstrumentStorageInterface,
	balanceStorage localStorage.BalanceStorageInterface,
	kafkaStorage localStorage.KafkaStorageInterface,
	userRepository *repositories.UserRepository,
	logger *log.Logger,
) *AppUseCase {
	return &AppUseCase{
		instrumentStorage: instrumentStorage,
		balanceStorage:    balanceStorage,
		kafkaStorage:      kafkaStorage,
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
	// Проверка на существование: если уже есть баланс, не обнуляем и возвращаем ошибку
	if u.balanceStorage.GetInstrumentBalance(instrument.InstrumentId).IsPositive() ||
		!u.balanceStorage.GetInstrumentBalance(instrument.InstrumentId).IsZero() {
		return uuid.Nil, fmt.Errorf("Баланс создан для инструмента с идентификатором")
	}

	// Сохраняем инструмент
	u.instrumentStorage.SaveInstrument(instrument)
	// Инициализируем баланс
	u.balanceStorage.SaveInstrumentBalance(instrument.InstrumentId)
	// Начисляем 100_000_000 рублей одним сложением
	u.balanceStorage.AddInstrumentBalance(instrument.InstrumentId, decimal.NewFromInt(100_000_000))

	// Возвращаем тот же идентификатор
	return instrument.InstrumentId, nil
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
