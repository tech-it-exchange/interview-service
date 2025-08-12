package localStorage

import (
	"github.com/google/uuid"
	"interview-service/port/outbound/localStorage"
	"sync"
)

type WorkerStorage struct {
	activeSpotInstrumentsMap   map[uuid.UUID]bool
	needSpotInstrumentDispatch bool
	mu                         sync.RWMutex
}

func NewWorkerStorage() localStorage.WorkerStorageInterface {
	return &WorkerStorage{
		activeSpotInstrumentsMap: make(map[uuid.UUID]bool),
	}
}

func (s *WorkerStorage) HasActiveSpotInstrument(instrumentId uuid.UUID) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.activeSpotInstrumentsMap[instrumentId]

	return ok
}

func (s *WorkerStorage) SaveActiveSpotInstrument(instrumentId uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.activeSpotInstrumentsMap[instrumentId] = true
}

func (s *WorkerStorage) SetSpotNeedInstrumentDispatch(value bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.needSpotInstrumentDispatch = value
}

func (s *WorkerStorage) IsNeedSpotInstrumentDispatch() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.needSpotInstrumentDispatch
}
