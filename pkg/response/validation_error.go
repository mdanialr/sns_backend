package response

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// validationError standard object that hold error from validation.
type validationError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// NewValidationErrors return ready to display error validations.
func NewValidationErrors(validations validator.ValidationErrors) []validationError {
	var validErr []validationError
	for _, valid := range validations {
		validErr = append(validErr, validationError{
			Name:    strings.ToLower(valid.Field()),
			Message: errMsgMapping(valid),
		})
	}
	return validErr
}

// errMsgMapping custom error message constructor from validator.
func errMsgMapping(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "required"
	case "len":
		return "length should be " + fe.Param()
	case "numeric":
		return "should be numeric"
	case "max":
		return "should be equal or less then " + fe.Param()
	case "min":
		return "should be equal or more then " + fe.Param()
	}
	return fe.Error()
}
