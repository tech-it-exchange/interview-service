package localStorage

import "github.com/google/uuid"

type WorkerStorageInterface interface {
	HasActiveSpotInstrument(instrumentId uuid.UUID) bool
	SaveActiveSpotInstrument(instrumentId uuid.UUID)
	SetSpotNeedInstrumentDispatch(value bool)
	IsNeedSpotInstrumentDispatch() bool
}
