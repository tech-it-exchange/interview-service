package events

import (
	"context"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"sync"
)

func (w *EventsWorker) balanceHandler(appCtx context.Context, instrumentId uuid.UUID) {
	var wgGroup sync.WaitGroup
	wgGroup.Add(1)

	ch := w.workerChannelManager.GetBalanceChannel(instrumentId)
	w.workerChannelManager.TriggerBalance(instrumentId)

	go func() {
		defer wgGroup.Done()

		for {
			select {
			case <-appCtx.Done():
				return
			case <-ch:
				// симуляция работы поступающего баланса из вне
				w.logger.Info("Добавление 100 000 000 рублей")
				for i := 1; i <= 100_000_000; i++ {
					go w.handleBalance(instrumentId, 1)
				}
				w.logger.Info("Добавлено 100 000 000 рублей инструменту", instrumentId)
			}
		}
	}()

	wgGroup.Wait()
}

func (w *EventsWorker) handleBalance(instrumentId uuid.UUID, balance int64) {
	rubDecimal := decimal.NewFromInt(balance)
	w.balanceStorage.AddInstrumentBalance(instrumentId, rubDecimal)
}
