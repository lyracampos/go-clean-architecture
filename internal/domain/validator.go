package domain

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var _ Validator = (*validatorService)(nil)

type Validator interface {
	Validate(input interface{}) error
}

type validatorService struct {
	validate *validator.Validate
}

func NewValidatorService() *validatorService {
	validate := validator.New()

	return &validatorService{
		validate: validate,
	}
}

func (v *validatorService) Validate(input interface{}) error {
	var errorList []FieldError
	if err := v.validate.Struct(input); err != nil {
		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			errorList = toFieldError(valErrs, errorList)
		}

		return NewValidationError(errorList)

	}
	return nil
}

func toFieldError(valErrs validator.ValidationErrors, errors []FieldError) []FieldError {
	for _, err := range valErrs {
		errors = append(errors, NewFieldError(err))
	}

	return errors
}
