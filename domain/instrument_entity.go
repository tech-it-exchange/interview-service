package domain

import (
	"github.com/google/uuid"
)

type InstrumentEntity struct {
	InstrumentId uuid.UUID
	IsListed     bool
	IsActive     bool
}
