// Package validator provides request validation utilities.
package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate checks a struct and returns a user-friendly error message.
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var valErr validator.ValidationErrors
	if !errors.As(err, &valErr) {
		return fmt.Errorf("validation error: %w", err)
	}

	var messages []string
	for _, e := range valErr {
		messages = append(messages, formatError(e))
	}
	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

func formatError(e validator.FieldError) string {
	field := e.Field()
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("field '%s' is required", field)
	case "email":
		return fmt.Sprintf("field '%s' must be a valid email", field)
	case "min":
		return fmt.Sprintf("field '%s' must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("field '%s' must be at most %s characters", field, e.Param())
	case "url":
		return fmt.Sprintf("field '%s' must be a valid URL", field)
	default:
		return fmt.Sprintf("field '%s' failed validation on '%s'", field, e.Tag())
	}
}
