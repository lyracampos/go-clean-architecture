package domain

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	ErrUserDoesNotExist  = errors.New("user does not exist")
	ErrEmailAlreadyInUse = errors.New("email is arealdy in use")
)

// ValidationError is thrown when a validation error happened.
type ValidationError struct {
	validationErrors []FieldError
}

// FieldError maps each field and validation error.
type FieldError struct {
	validator.FieldError
}

// NewValidationError instantiates a new validation error.
func NewValidationError(errors []FieldError) error {
	return &ValidationError{
		validationErrors: errors,
	}
}

// Error is the implementation of Error interface.
func (v *ValidationError) Error() string {
	fieldErrors := ""
	for _, validationError := range v.validationErrors {
		fieldErrors += validationError.Error() + "; "
	}

	return fieldErrors
}

func (v *ValidationError) ValidationErrors() []FieldError {
	return v.validationErrors
}

func NewFieldError(fieldError validator.FieldError) FieldError {
	return FieldError{
		FieldError: fieldError,
	}
}

func (f *FieldError) Error() string {
	switch f.FieldError.Tag() {
	case "required":
		return fmt.Sprintf("the field '%s' should not be empty", f.FieldError.Field())
	case "min", "max", "maxBytes":
		return fmt.Sprintf("the field '%s' %s size is %s", f.FieldError.Field(), f.FieldError.Tag(), f.FieldError.Param())
	case "oneof":
		return fmt.Sprintf("the field '%s' is not valid, expected one of [%s]", f.FieldError.Field(), f.FieldError.Param())
	case "uuid":
		return fmt.Sprintf("the field '%s' is not a valid UUID", f.FieldError.Field())
	default:
		return fmt.Sprintf("the field '%s' is invalid", f.FieldError.Field())
	}
}
