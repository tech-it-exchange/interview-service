package dto

import (
	"github.com/google/uuid"
)

type CreateInstrumentRequestDto struct {
	InstrumentId uuid.UUID `json:"instrumentId" validate:"required" example:"0fae8035-7e4b-4f60-b1da-e3752d53f2f4"`
	IsActive     bool      `json:"isActive" example:"true"`
	IsListed     bool      `json:"isListed" example:"false"`
}
