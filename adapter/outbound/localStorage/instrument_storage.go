package localStorage

import (
	"github.com/google/uuid"
	"interview-service/domain"
	"interview-service/port/outbound/localStorage"
	"sync"
	"sync/atomic"
)

type instrumentStorage struct {
	mu                 sync.RWMutex
	spotInstrumentsMap atomic.Value
}

func NewInstrumentStorage() localStorage.InstrumentStorageInterface {
	storage := &instrumentStorage{}
	storage.spotInstrumentsMap.Store(make(map[string]*domain.InstrumentEntity))

	return storage
}

func (s *instrumentStorage) GetSpotInstrumentsMap() map[string]*domain.InstrumentEntity {
	return s.spotInstrumentsMap.Load().(map[string]*domain.InstrumentEntity)
}

func (s *instrumentStorage) GetActiveSpotInstrumentIds() []uuid.UUID {
	current := s.spotInstrumentsMap.Load().(map[string]*domain.InstrumentEntity)

	var ids []uuid.UUID
	for _, instrument := range current {
		if instrument == nil || !instrument.IsActive || !instrument.IsListed {
			continue
		}

		ids = append(ids, instrument.InstrumentId)
	}

	return ids
}

// SaveInstrument сохраняет или обновляет инструмент в хранилище
func (s *instrumentStorage) SaveInstrument(instrument *domain.InstrumentEntity) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current := s.spotInstrumentsMap.Load().(map[string]*domain.InstrumentEntity)

	newMap := make(map[string]*domain.InstrumentEntity, len(current)+1)
	for k, v := range current {
		newMap[k] = v
	}

	newMap[instrument.InstrumentId.String()] = instrument

	s.spotInstrumentsMap.Store(newMap)
}
