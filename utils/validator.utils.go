package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type FieldError struct {
	Field   string `json:"field"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message"`
}

func ValidateStruct(s any) []FieldError {
	if err := validate.Struct(s); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return translateValidationErrors(ve)
		}
		return []FieldError{{
			Field:   "_",
			Message: "Invalid input payload",
		}}
	}
	return nil
}

func translateValidationErrors(ve validator.ValidationErrors) []FieldError {
	out := make([]FieldError, 0, len(ve))
	for _, fe := range ve {
		out = append(out, FieldError{
			Field:   fe.Field(),
			Param:   fe.Param(),
			Message: messageFor(fe),
		})
	}
	return out
}

func messageFor(fe validator.FieldError) string {
	field := fe.Field()
	param := fe.Param()

	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "url":
		return field + " must be a valid URL"
	case "uuid", "uuid4":
		return field + " must be a valid UUID"
	case "alpha":
		return field + " must contain only letters"
	case "alphanum":
		return field + " must contain only letters and numbers"
	case "numeric":
		return field + " must be numeric"

	case "min":
		return field + " must be at least " + param
	case "max":
		return field + " must be at most " + param
	case "len":
		return field + " must be exactly " + param

	case "gte":
		return field + " must be greater than or equal to " + param
	case "lte":
		return field + " must be less than or equal to " + param
	case "gt":
		return field + " must be greater than " + param
	case "lt":
		return field + " must be less than " + param

	case "oneof":
		opts := strings.ReplaceAll(param, " ", ", ")
		return field + " must be one of: " + opts

	case "eqfield":
		return field + " must be equal to " + param
	case "nefield":
		return field + " must be different from " + param
	}

	return field + " failed on the '" + fe.Tag() + "' rule"
}
