package routes

import (
	"github.com/gofiber/fiber/v2"
	api "github.com/mdanialr/sns_backend/internal/api/v1/send"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/mdanialr/sns_backend/internal/middleware"
	"github.com/mdanialr/sns_backend/internal/service"
)

// SendRoutes register all routes that handle Send data.
func SendRoutes(r fiber.Router, conf *service.Config, db *database.Queries) {
	r.Post("/",
		middleware.CreateSendValidation,
		api.CreateSend(conf, db),
	)
}
