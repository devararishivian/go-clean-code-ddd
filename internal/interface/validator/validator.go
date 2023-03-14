package validator

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

type Validator struct {
	validator *validator.Validate
}

type ErrorResponse struct {
	Message         string            `json:"message"`
	ValidationError []validationError `json:"validationError"`
}

type validationError struct {
	FailedField string `json:"field"`
	Tag         string `json:"validationTag"`
	Value       string `json:"value"`
}

func NewValidator() *Validator {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("json")
	})

	return &Validator{validator: validate}
}

func (cv *Validator) Validate(i interface{}) ErrorResponse {
	var errors ErrorResponse

	if err := cv.validator.Struct(i); err != nil {
		for _, fieldError := range err.(validator.ValidationErrors) {
			each := validationError{
				FailedField: fieldError.Field(),
				Tag:         fieldError.Tag(),
				Value:       fieldError.Param(),
			}

			errors.ValidationError = append(errors.ValidationError, each)
		}

		errors.Message = "error while validating input"
	}

	return errors
}
