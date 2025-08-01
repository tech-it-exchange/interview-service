package controllers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"interview-service/adapter/inbound/fasthttp/dto"
	httpMappers "interview-service/adapter/inbound/fasthttp/mappers"
	"interview-service/application/usecase"
	fasthttp2 "interview-service/infrastructure/fasthttp"
)

type ClientOrderController struct {
	validate   *validator.Validate
	appUseCase *usecase.AppUseCase
	httpMapper *httpMappers.HttpMapper
}

func NewClientOrderController(
	validate *validator.Validate,
	workerUseCase *usecase.AppUseCase,
	httpMapper *httpMappers.HttpMapper,
) *ClientOrderController {
	return &ClientOrderController{
		validate:   validate,
		appUseCase: workerUseCase,
		httpMapper: httpMapper,
	}
}

// GetLoss @Summary Get loss
// @Description Get loss.
// @Tags client
// @Security BearerToken
// @Success 200 {object} dto.LossResponse
// @Router /api/v2/client/user/get-loss [GET]
func (c *ClientOrderController) GetLoss(ctx *fasthttp.RequestCtx) {
	result, err := c.appUseCase.GetLoss(ctx)
	if err != nil {
		fasthttp2.WriteErrorResponse(
			ctx,
			err.Error(),
			fasthttp.StatusOK,
			err,
		)

		return
	}

	fasthttp2.WriteResponse(ctx, dto.LossResponse{
		Loss: result,
	})
}

// CreateOrder @Summary Create order
// @Description Create order.
// @Tags client
// @Security BearerToken
// @Param  request body dto.CreateOrderRequestDto true "body"
// @Success 200 {object} dto.CreateOrderResponseDto
// @Router /api/v2/client/order/create [post]
func (c *ClientOrderController) CreateOrder(ctx *fasthttp.RequestCtx) {
	request := &dto.CreateOrderRequestDto{}

	err := json.Unmarshal(ctx.PostBody(), &request)
	if err != nil {
		fasthttp2.WriteErrorResponse(
			ctx,
			err.Error(),
			fasthttp.StatusBadRequest,
			err,
		)

		return
	}

	err = c.validate.Struct(request)
	if err != nil {
		fasthttp2.WriteErrorResponse(
			ctx,
			err.Error(),
			fasthttp.StatusBadRequest,
			err,
		)

		return
	}

	orderId, err := c.appUseCase.CreateOrder(ctx, request)
	if err != nil {
		fasthttp2.WriteErrorResponse(
			ctx,
			err.Error(),
			fasthttp.StatusInternalServerError,
			err,
		)

		return
	}

	fasthttp2.WriteResponse(ctx, c.httpMapper.MapAppCreateOrderToResponseDto(orderId))
}

// CreateInstrument @Summary Create instrument
// @Description Create instrument.
// @Tags client
// @Security BearerToken
// @Param  request body dto.CreateInstrumentRequestDto true "body"
// @Success 200 {object} dto.CreateInstrumentResponseDto
// @Router /api/v2/client/instrument/create [post]
func (c *ClientOrderController) CreateInstrument(ctx *fasthttp.RequestCtx) {
	request := &dto.CreateInstrumentRequestDto{}

	err := json.Unmarshal(ctx.PostBody(), &request)
	if err != nil {
		fasthttp2.WriteErrorResponse(
			ctx,
			err.Error(),
			fasthttp.StatusBadRequest,
			err,
		)

		return
	}

	err = c.validate.Struct(request)
	if err != nil {
		fasthttp2.WriteErrorResponse(
			ctx,
			err.Error(),
			fasthttp.StatusBadRequest,
			err,
		)

		return
	}

	instrumentId, err := c.appUseCase.CreateInstrument(ctx, c.httpMapper.MapCreateInstrumentRequestToDomain(request))
	if err != nil {
		fasthttp2.WriteErrorResponse(
			ctx,
			err.Error(),
			fasthttp.StatusInternalServerError,
			err,
		)

		return
	}

	fasthttp2.WriteResponse(ctx, c.httpMapper.MapAppCreateInstrumentToResponseDto(instrumentId))
}

// GetUserName @Summary Get username
// @Description Get username.
// @Tags client
// @Security BearerToken
// @Success 200 {object} dto.UserResponse
// @Router /api/v2/client/user/by/{name} [GET]
func (c *ClientOrderController) GetUserName(ctx *fasthttp.RequestCtx) {
	name, _ := ctx.UserValue("name").(string)
	result := c.appUseCase.GetUserName(ctx, name)

	fasthttp2.WriteResponse(ctx, &dto.UserResponse{
		Name: result,
	})
}

// GetInstrumentBalance @Summary Get instrument balance
// @Description Get instrument balance.
// @Tags client
// @Security BearerToken
// @Success 200 {object} dto.InstrumentResponse
// @Router /api/v2/client/instrument/balance/{instrumentId} [GET]
func (c *ClientOrderController) GetInstrumentBalance(ctx *fasthttp.RequestCtx) {
	instrumentId, _ := ctx.UserValue("instrumentId").(string)
	result := c.appUseCase.GetInstrumentBalance(ctx, uuid.MustParse(instrumentId))

	fasthttp2.WriteResponse(ctx, dto.InstrumentResponse{
		Balance: result,
	})
}

// GetActiveConsumerTopics @Summary Get active consumer topics
// @Description Get active consumer topics.
// @Tags client
// @Security BearerToken
// @Success 200 {object} []string
// @Router /api/v2/client/kafka/topics [GET]
func (c *ClientOrderController) GetActiveConsumerTopics(ctx *fasthttp.RequestCtx) {
	result := c.appUseCase.GetActiveConsumerTopics()

	fasthttp2.WriteResponse(ctx, result)
}
