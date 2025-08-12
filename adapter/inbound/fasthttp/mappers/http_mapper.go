package mappers

import (
	"github.com/google/uuid"
	"interview-service/adapter/inbound/fasthttp/dto"
	"interview-service/domain"
)

type HttpMapper struct{}

func NewHttpMapper() *HttpMapper {
	return &HttpMapper{}
}

func (m *HttpMapper) MapAppCreateOrderToResponseDto(
	orderId uuid.UUID,
) *dto.CreateOrderResponseDto {
	return &dto.CreateOrderResponseDto{
		OrderId: orderId.String(),
	}
}

func (m *HttpMapper) MapCreateInstrumentRequestToDomain(
	requestDto *dto.CreateInstrumentRequestDto,
) *domain.InstrumentEntity {
	return &domain.InstrumentEntity{
		InstrumentId: requestDto.InstrumentId,
		IsListed:     requestDto.IsListed,
		IsActive:     requestDto.IsActive,
	}
}

func (m *HttpMapper) MapAppCreateInstrumentToResponseDto(
	instrumentId uuid.UUID,
) *dto.CreateInstrumentResponseDto {
	return &dto.CreateInstrumentResponseDto{
		InstrumentId: instrumentId.String(),
	}
}
