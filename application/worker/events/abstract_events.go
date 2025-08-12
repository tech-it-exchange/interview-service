package events

import (
	"github.com/google/uuid"
)

// triggerEvent Общий метод для триггера по необработанному количеству ивентов
func triggerEvent(
	countInboxEvents int,
	ch chan struct{},
	concurrent int,
	instrumentId uuid.UUID,
	trigger func(instrumentId uuid.UUID),
) {
	if countInboxEvents == 0 {
		return
	}

	pending := len(ch) * concurrent
	count := countInboxEvents - pending

	if count <= 0 {
		return
	}

	availableSlots := cap(ch) - len(ch)
	toTrigger := min(count, availableSlots)

	for i := 0; i < toTrigger; i++ {
		trigger(instrumentId)
	}
}
