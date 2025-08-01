package localStorage

import (
	"github.com/google/uuid"
	"interview-service/domain"
)

type InstrumentStorageInterface interface {
	GetSpotInstrumentsMap() map[string]*domain.InstrumentEntity
	GetActiveSpotInstrumentIds() []uuid.UUID
	SaveInstrument(instrument *domain.InstrumentEntity)
}
