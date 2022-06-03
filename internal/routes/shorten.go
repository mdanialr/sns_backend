package routes

import (
	"github.com/gofiber/fiber/v2"
	api "github.com/mdanialr/sns_backend/internal/api/v1/shorten"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/mdanialr/sns_backend/internal/middleware"
)

// ShortenRoutes register all routes that handle Shorten data.
func ShortenRoutes(r fiber.Router, db *database.Queries) {
	r.Get("/",
		api.ListShorten(db),
	)
	r.Get("/:id",
		api.GetShortenDetail(db),
	)
	r.Post("/",
		middleware.CreateShortenValidation,
		api.CreateShorten(db),
	)
	r.Patch("/:id",
		middleware.IsIdExistsValidation(db),
		middleware.CreateShortenValidation,
		api.UpdateShorten(db),
	)
	r.Put("/:id",
		middleware.IsIdExistsValidation(db),
		middleware.CreateShortenValidation,
		api.UpdateShorten(db),
	)
	r.Delete("/:id",
		middleware.IsIdExistsValidation(db),
		api.DeleteShorten(db),
	)
}
