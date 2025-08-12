package dto

import (
	"github.com/google/uuid"
)

type CreateOrderRequestDto struct {
	InstrumentId uuid.UUID `json:"instrumentId" validate:"required" example:"0fae8035-7e4b-4f60-b1da-e3752d53f2f4"`
	// ask или bid
	Side  string `json:"side" validate:"required,side" example:"ask"`
	Qty   string `json:"qty" example:"1"`
	Price string `json:"price" example:"100"`
}
