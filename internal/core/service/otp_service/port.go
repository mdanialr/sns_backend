package otp_service

import (
	"context"

	req "github.com/mdanialr/sns_backend/internal/requests"
)

// IService an interface that should be used when dealing with otp.
type IService interface {
	// ValidateOTP validate the given request that should have otp code inside.
	// Should always return false either for invalid token or any error in
	// order to verify the otp code.
	ValidateOTP(context.Context, *req.OTP) bool
	// GetJWT return the string encoded jwt token that's ready to be served to
	// response, also return error if any.
	GetJWT(context.Context) (string, error)
}
