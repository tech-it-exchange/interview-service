package events

import (
	"context"
	"log"

	"interview-service/application/channels"
	"interview-service/application/worker/abstract"
	"interview-service/port/outbound/localStorage"

	"github.com/google/uuid"
)

type EventsWorker struct {
	workerChannelManager     *channels.WorkerChannelManager
	workerStorage            localStorage.WorkerStorageInterface
	instrumentStorage        localStorage.InstrumentStorageInterface
	instrumentChannelManager *channels.InstrumentChannelManager
	balanceStorage           localStorage.BalanceStorageInterface
	logger                   *log.Logger
}

func NewEventsWorker(
	workerChannelManager *channels.WorkerChannelManager,
	workerStorage localStorage.WorkerStorageInterface,
	instrumentStorage localStorage.InstrumentStorageInterface,
	instrumentChannelManager *channels.InstrumentChannelManager,
	balanceStorage localStorage.BalanceStorageInterface,
	logger *log.Logger,
) abstract.WorkerInterface {
	return &EventsWorker{
		workerChannelManager:     workerChannelManager,
		workerStorage:            workerStorage,
		instrumentStorage:        instrumentStorage,
		instrumentChannelManager: instrumentChannelManager,
		balanceStorage:           balanceStorage,
		logger:                   logger,
	}
}

// Start Запускает воркеры
func (w *EventsWorker) Start(appCtx context.Context) {
	instrumentIds := w.instrumentStorage.GetActiveSpotInstrumentIds()

	for _, instrumentId := range instrumentIds {
		w.start(appCtx, instrumentId)
		w.workerStorage.SaveActiveSpotInstrument(instrumentId)
	}

	w.workerStorage.SetSpotNeedInstrumentDispatch(true)

	go w.listenForNewSpotInstrument(appCtx)
}

// ListenForNewSpotInstrument Слушает спот инструменты и запускает новые потребители по мере их поступления
func (w *EventsWorker) listenForNewSpotInstrument(appCtx context.Context) {
	for {
		select {
		case <-appCtx.Done():
			w.logger.Printf("Остановка ListenForNewSpotInstrument: %v", appCtx.Err())
			return
		case instrumentId, ok := <-w.instrumentChannelManager.GetSpotSyncEventInstrumentIdChannel():
			if !ok {
				w.logger.Printf("Пропуск, канал новых спот инструментов закрыт")
				continue
			}

			if w.workerStorage.HasActiveSpotInstrument(instrumentId) {
				w.logger.Printf("Обработчики ивентов для спот инструмента уже существуют: %s", instrumentId)

				continue
			}

			w.start(appCtx, instrumentId)
			w.workerStorage.SaveActiveSpotInstrument(instrumentId)

			w.logger.Printf("Запущены обработчики ивентов для спот инструмента: %s", instrumentId)
		}
	}
}

// start Запускает обработчики
func (w *EventsWorker) start(appCtx context.Context, instrumentId uuid.UUID) {
	go w.balanceHandler(appCtx, instrumentId)
}
