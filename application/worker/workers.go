package worker

import (
	"context"
	"interview-service/adapter/outbound/localStorage"
	"interview-service/application/channels"
	"interview-service/application/service"
	"interview-service/application/worker/abstract"
	"interview-service/application/worker/cache"
	"interview-service/application/worker/events"
	"log/slog"
	"time"
)

type Workers struct {
	instrumentCache abstract.WorkerInterface
	syncEvents      abstract.WorkerInterface
}

func NewWorkers(
	storageContainer *localStorage.StorageContainer,
	services *service.Services,
	logger *slog.Logger,
	refreshInstrumentsTimeout time.Duration,
	workerChannelManager *channels.WorkerChannelManager,
	instrumentChannelManager *channels.InstrumentChannelManager,
) WorkersInterface {
	instrumentCacheWorker := cache.NewInstrumentCacheWorker(
		storageContainer.Instrument,
		refreshInstrumentsTimeout,
		instrumentChannelManager,
		services.Kafka,
		storageContainer.Worker,
		logger,
	)

	syncEvents := events.NewEventsWorker(
		workerChannelManager,
		storageContainer.Worker,
		storageContainer.Instrument,
		instrumentChannelManager,
		storageContainer.Balance,
		logger,
	)

	return &Workers{
		instrumentCache: instrumentCacheWorker,
		syncEvents:      syncEvents,
	}
}

func (w *Workers) StartWorkers(ctx context.Context) {
	w.instrumentCache.Start(ctx)

	w.syncEvents.Start(ctx)
}
