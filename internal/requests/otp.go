package requests

import (
	"github.com/go-playground/validator/v10"
)

// OTP request that should be used inside endpoint `otp`
type OTP struct {
	Code string `json:"code" validate:"required,numeric,len=6"`
}

// Validate do validate the struct that should be parsed from request body.
func (o *OTP) Validate() validator.ValidationErrors {
	if err := validate.Struct(o); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}
