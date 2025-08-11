package localStorage

import (
	"interview-service/port/outbound/localStorage"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type balanceStorage struct {
	balanceInstrumentsMap map[uuid.UUID]decimal.Decimal
}

func NewBalanceStorage() localStorage.BalanceStorageInterface {
	return &balanceStorage{
		balanceInstrumentsMap: make(map[uuid.UUID]decimal.Decimal),
	}
}

func (s *balanceStorage) SaveInstrumentBalance(instrumentId uuid.UUID) {
	// Инициализируем баланс нулем, если не существует
	if _, ok := s.balanceInstrumentsMap[instrumentId]; !ok {
		s.balanceInstrumentsMap[instrumentId] = decimal.Zero
	}
}

func (s *balanceStorage) GetInstrumentBalance(instrumentId uuid.UUID) decimal.Decimal {
	return s.balanceInstrumentsMap[instrumentId]
}

func (s *balanceStorage) AddInstrumentBalance(instrumentId uuid.UUID, balance decimal.Decimal) {
	current := s.balanceInstrumentsMap[instrumentId]
	s.balanceInstrumentsMap[instrumentId] = current.Add(balance)
}
