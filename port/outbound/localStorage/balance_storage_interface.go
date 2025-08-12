package localStorage

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BalanceStorageInterface interface {
	SaveInstrumentBalance(instrumentId uuid.UUID)
	GetInstrumentBalance(instrumentId uuid.UUID) decimal.Decimal
	AddInstrumentBalance(instrumentId uuid.UUID, balance decimal.Decimal)
}
