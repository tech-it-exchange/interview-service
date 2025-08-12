package events

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"interview-service/application/channels"
	"interview-service/application/worker/abstract"
	"interview-service/port/outbound/localStorage"
)

type EventsWorker struct {
	workerChannelManager     *channels.WorkerChannelManager
	workerStorage            localStorage.WorkerStorageInterface
	instrumentStorage        localStorage.InstrumentStorageInterface
	instrumentChannelManager *channels.InstrumentChannelManager
	balanceStorage           localStorage.BalanceStorageInterface
	logger                   *slog.Logger
}

func NewEventsWorker(
	workerChannelManager *channels.WorkerChannelManager,
	workerStorage localStorage.WorkerStorageInterface,
	instrumentStorage localStorage.InstrumentStorageInterface,
	instrumentChannelManager *channels.InstrumentChannelManager,
	balanceStorage localStorage.BalanceStorageInterface,
	logger *slog.Logger,
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
			w.logger.Info("Остановка ListenForNewSpotInstrument:", appCtx.Err())
			return
		case instrumentId, ok := <-w.instrumentChannelManager.GetSpotSyncEventInstrumentIdChannel():
			if !ok {
				w.logger.Error("Пропуск, канал новых спот инструментов закрыт")
				continue
			}

			if w.workerStorage.HasActiveSpotInstrument(instrumentId) {
				w.logger.Info(
					"Обработчики ивентов для спот инструмента уже существуют",
					"instrumentId",
					instrumentId,
				)

				continue
			}

			w.start(appCtx, instrumentId)
			w.workerStorage.SaveActiveSpotInstrument(instrumentId)

			w.logger.Info("Запущены обработчики ивентов для спот инструмента", "instrumentId", instrumentId)
		}
	}
}

// start Запускает обработчики
func (w *EventsWorker) start(appCtx context.Context, instrumentId uuid.UUID) {
	go w.balanceHandler(appCtx, instrumentId)
}
