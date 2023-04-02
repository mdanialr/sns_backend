package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/sns_backend/internal/app/adapter/http/auth_handler"
	"github.com/mdanialr/sns_backend/internal/app/adapter/http/shorten_handler"
	"github.com/mdanialr/sns_backend/internal/core/repository/otp_repository"
	"github.com/mdanialr/sns_backend/internal/core/service/otp_service"
	"github.com/mdanialr/sns_backend/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// HttpHandlers handlers that use HTTP as the controller/handler.
type HttpHandlers struct {
	R      fiber.Router
	DB     *gorm.DB
	Config *viper.Viper
	Log    logger.Writer
}

func (h *HttpHandlers) SetupRouter() {
	h.R.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("OK")
	})
	h.R.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("SNS Backend API")
	})
	// currently use v1
	apiV1 := h.R.Group("/v1")

	// init repositories
	otpRepo := otp_repository.NewOTPRepository(h.DB)

	// init services
	otpSvc := otp_service.NewOTPService(h.Config, h.Log, otpRepo)

	// init handlers
	auth_handler.NewAuthHandler(apiV1, otpSvc)         // /auth/*
	shorten_handler.NewShortenHandler(apiV1, h.Config) // /shorten/*
}
