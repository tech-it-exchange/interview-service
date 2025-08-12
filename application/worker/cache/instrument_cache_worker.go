package cache

import (
	"context"
	"interview-service/application/channels"
	"interview-service/application/service"
	"interview-service/application/worker/abstract"
	"interview-service/domain"
	localStoragePort "interview-service/port/outbound/localStorage"
	"log/slog"
	"time"
)

type InstrumentCacheWorker struct {
	instrumentStorage         localStoragePort.InstrumentStorageInterface
	refreshInstrumentsTimeout time.Duration
	instrumentChannelManager  *channels.InstrumentChannelManager
	kafkaService              *service.KafkaService
	workerStorage             localStoragePort.WorkerStorageInterface
	logger                    *slog.Logger
}

func NewInstrumentCacheWorker(
	instrumentStorage localStoragePort.InstrumentStorageInterface,
	refreshInstrumentsTimeout time.Duration,
	instrumentChannelManager *channels.InstrumentChannelManager,
	kafkaService *service.KafkaService,
	workerStorage localStoragePort.WorkerStorageInterface,
	logger *slog.Logger,
) abstract.WorkerInterface {
	return &InstrumentCacheWorker{
		instrumentStorage:         instrumentStorage,
		refreshInstrumentsTimeout: refreshInstrumentsTimeout,
		instrumentChannelManager:  instrumentChannelManager,
		kafkaService:              kafkaService,
		workerStorage:             workerStorage,
		logger:                    logger,
	}
}

// Start Запускает воркеры
func (w *InstrumentCacheWorker) Start(appCtx context.Context) {
	w.getAndUpdateInstrumentsData()

	go w.updateInstrumentsData(appCtx)
}

// updateInstrumentsData Обновление данных в локальном хранилище инструментов, данными из redis
func (w *InstrumentCacheWorker) updateInstrumentsData(appCtx context.Context) {
	ticker := time.NewTicker(w.refreshInstrumentsTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-appCtx.Done():
			w.logger.Info("Контекст отменён, остановка обновления инструментов из redis")

			return
		case <-ticker.C:
			w.getAndUpdateInstrumentsData()
		}
	}
}

// getAndUpdateInstrumentsData Получает и обновляет инструменты
func (w *InstrumentCacheWorker) getAndUpdateInstrumentsData() {
	spotInstrumentsMap := w.instrumentStorage.GetSpotInstrumentsMap()

	w.filterAndDispatchSpotInstrumentUuids(spotInstrumentsMap)
}

// filterAndDispatchSpotInstruments Фильтрует и отправляет идентификаторы спот инструментов в канал
func (w *InstrumentCacheWorker) filterAndDispatchSpotInstrumentUuids(spotInstrumentsMap map[string]*domain.InstrumentEntity) {
	for _, spotInstrument := range spotInstrumentsMap {
		if !spotInstrument.IsListed || !spotInstrument.IsActive {
			continue
		}

		instrumentId := spotInstrument.InstrumentId

		if w.kafkaService.IsNeedSpotInstrumentDispatch() && !w.kafkaService.HasActiveSpotInstrument(instrumentId) {
			w.instrumentChannelManager.SendKafkaSpotInstrumentId(instrumentId)
		}

		if w.workerStorage.IsNeedSpotInstrumentDispatch() && !w.workerStorage.HasActiveSpotInstrument(instrumentId) {
			w.instrumentChannelManager.SendSpotSyncEventSpotInstrumentId(instrumentId)
		}
	}
}
