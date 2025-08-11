package events

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (w *EventsWorker) balanceHandler(appCtx context.Context, instrumentId uuid.UUID) {
	ch := w.workerChannelManager.GetBalanceChannel(instrumentId)
	w.workerChannelManager.TriggerBalance(instrumentId)

	for {
		select {
		case <-appCtx.Done():
			return
		case <-ch:
			// Начисляем сумму одним действием, если баланс ещё не начислен
			if w.balanceStorage.GetInstrumentBalance(instrumentId).IsZero() {
				w.handleBalance(instrumentId, 100_000_000)
				w.logger.Printf("Добавлено 100 000 000 рублей инструменту %s", instrumentId)
			}
		}
	}
}

func (w *EventsWorker) handleBalance(instrumentId uuid.UUID, balance int64) {
	rubDecimal := decimal.NewFromInt(balance)
	w.balanceStorage.AddInstrumentBalance(instrumentId, rubDecimal)
}
