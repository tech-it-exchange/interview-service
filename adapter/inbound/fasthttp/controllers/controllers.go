package controllers

import (
	"github.com/go-playground/validator/v10"
	httpMappers "interview-service/adapter/inbound/fasthttp/mappers"
	httpValidators "interview-service/adapter/inbound/fasthttp/validators"
	"interview-service/application/usecase"
)

type Controllers struct {
	ClientOrderController *ClientOrderController
}

type Dependencies struct {
	UseCases *usecase.UseCases
}

func NewControllers(deps *Dependencies) (*Controllers, error) {
	mappers := httpMappers.NewMappers()
	requestValidator := validator.New()

	requestValidators := httpValidators.NewValidators(requestValidator)
	err := requestValidators.SpotCustomValidator.RegisterCustomValidators()
	if err != nil {
		return nil, err
	}

	clientOrderController := NewClientOrderController(
		requestValidator,
		deps.UseCases.Worker,
		mappers.Http,
	)

	return &Controllers{
		ClientOrderController: clientOrderController,
	}, nil
}
