package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/sns_backend/internal/app/adapter/http/auth_handler"
	"github.com/mdanialr/sns_backend/internal/app/adapter/http/shorten_handler"
	"github.com/mdanialr/sns_backend/internal/core/repository/otp_repository"
	"github.com/mdanialr/sns_backend/internal/core/repository/sns_repository"
	"github.com/mdanialr/sns_backend/internal/core/service/otp_service"
	"github.com/mdanialr/sns_backend/internal/core/service/shorten_service"
	"github.com/mdanialr/sns_backend/pkg/logger"
	"github.com/mdanialr/sns_backend/pkg/storage"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// HttpHandlers handlers that use HTTP as the controller/handler.
type HttpHandlers struct {
	R       fiber.Router
	DB      *gorm.DB
	Config  *viper.Viper
	Log     logger.Writer
	Storage storage.IStorage
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
	otpRepo := otp_repository.New(h.DB)
	snsRepo := sns_repository.New(h.DB)

	// init services
	otpSvc := otp_service.New(h.Config, h.Log, otpRepo)
	snsSvc := shorten_service.New(h.Log, snsRepo)

	// init handlers
	auth_handler.New(apiV1, otpSvc)              // /auth/*
	shorten_handler.New(apiV1, h.Config, snsSvc) // /shorten/*
}
