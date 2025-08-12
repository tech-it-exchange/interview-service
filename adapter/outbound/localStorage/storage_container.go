package localStorage

import (
	"interview-service/port/outbound/localStorage"
)

type StorageContainer struct {
	Instrument localStorage.InstrumentStorageInterface
	Kafka      localStorage.KafkaStorageInterface
	Worker     localStorage.WorkerStorageInterface
	Balance    localStorage.BalanceStorageInterface
}

func NewStorageContainer() *StorageContainer {
	instrument := NewInstrumentStorage()
	kafka := NewKafkaStorage()
	worker := NewWorkerStorage()
	balance := NewBalanceStorage()

	return &StorageContainer{
		Instrument: instrument,
		Kafka:      kafka,
		Worker:     worker,
		Balance:    balance,
	}
}
