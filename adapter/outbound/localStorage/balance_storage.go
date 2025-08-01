package localStorage

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"interview-service/port/outbound/localStorage"
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
	s.balanceInstrumentsMap[instrumentId] = decimal.Zero
}

func (s *balanceStorage) GetInstrumentBalance(instrumentId uuid.UUID) decimal.Decimal {
	return s.balanceInstrumentsMap[instrumentId]
}

func (s *balanceStorage) AddInstrumentBalance(instrumentId uuid.UUID, balance decimal.Decimal) {
	cacheBalance := s.balanceInstrumentsMap[instrumentId]

	cacheBalance.Add(balance)

	s.balanceInstrumentsMap[instrumentId] = cacheBalance
}
