package auth_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/sns_backend/internal/core/service/otp_service"
	req "github.com/mdanialr/sns_backend/internal/requests"
	"github.com/mdanialr/sns_backend/pkg/constant"
	resp "github.com/mdanialr/sns_backend/pkg/response"
)

type authHandler struct {
	route  fiber.Router
	otpSvc otp_service.IService
}

// New init all endpoints within `/auth`.
func New(route fiber.Router, otpSvc otp_service.IService) {
	otpH := &authHandler{route, otpSvc}

	api := otpH.route.Group("/auth")
	api.Post("/otp", otpH.Login)
}

// Login exchange given totp with a jwt.
func (a *authHandler) Login(c *fiber.Ctx) error {
	var ot = new(req.OTP)
	c.BodyParser(ot)
	// validate the request
	if err := ot.Validate(); err != nil {
		return resp.Error(c, resp.WithErrMsg(cons.InvalidPayload), resp.WithErrValidation(err))
	}
	// validate the incoming otp
	if !a.otpSvc.ValidateOTP(c.Context(), ot) {
		return resp.Error(c, resp.WithErrMsg(cons.InvalidOTP))
	}
	// get the signed jwt token
	token, err := a.otpSvc.GetJWT(c.Context())
	if err != nil {
		return resp.ErrorCode(c, fiber.StatusInternalServerError, resp.WithErr(err))
	}
	// construct the response
	tokenResp := map[string]string{
		"token": token,
	}

	return resp.Success(c, resp.WithData(tokenResp))
}
