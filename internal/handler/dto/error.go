package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Message string            `json:"message" example:"Invalid request body"`
	Errors  map[string]string `json:"errors,omitempty" example:"email: invalid format"`
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
	}
}

func NewValidationErrorResponse(err error) ErrorResponse {
	res := ErrorResponse{
		Message: "Validation failed",
		Errors:  make(map[string]string),
	}

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, f := range validationErrs {
			res.Errors[f.Field()] = fmt.Sprintf("failed on '%s' tag", f.Tag())
		}
	} else {
		res.Message = err.Error()
	}

	return res
}
