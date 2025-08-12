package validators

import (
	"github.com/go-playground/validator/v10"
	"interview-service/infrastructure/consts"
	"interview-service/infrastructure/types"
)

type CustomValidator struct {
	validate *validator.Validate
}

func NewCustomValidator(validate *validator.Validate) *CustomValidator {
	return &CustomValidator{
		validate: validate,
	}
}

func (v *CustomValidator) validateSide(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().Interface()

	side, ok := fieldValue.(string)
	if !ok {
		return false
	}

	validSides := []types.OrderSide{
		consts.Ask,
		consts.Bid,
	}

	for _, valid := range validSides {
		if types.OrderSide(side) != valid {
			continue
		}

		return true
	}

	return false
}

// RegisterCustomValidators Регистрация кастомных валидаторов
func (v *CustomValidator) RegisterCustomValidators() error {
	err := v.validate.RegisterValidation("side", v.validateSide)
	if err != nil {
		return err
	}

	return nil
}
