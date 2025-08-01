package validators

import (
	"github.com/go-playground/validator/v10"
)

type Validators struct {
	SpotCustomValidator *CustomValidator
}

func NewValidators(validate *validator.Validate) *Validators {
	spotCustomValidator := NewCustomValidator(validate)

	return &Validators{
		SpotCustomValidator: spotCustomValidator,
	}
}
