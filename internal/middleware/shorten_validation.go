package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/sns_backend/internal/api/v1/shorten"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/mdanialr/sns_backend/internal/service"
)

// CreateShortenValidation middleware that handle necessary validation for creating new shorten.
func CreateShortenValidation(c *fiber.Ctx) error {
	var payload api.CreateShortenPayload
	if err := c.BodyParser(&payload); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("failed to parse json payload: %s", err),
		})
	}

	// 1st validation: `target` should not be empty
	if payload.Target == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "`target` field should not be empty",
		})
	}

	// 2nd validation: if `url` field is not provided then substitute it with random string
	if payload.Url == "" {
		payload.Url = service.RandomString(8)
	}

	return c.Next()
}

// IsIdExistsValidation middleware that handle validation that make sure the ID in URL is exists in the database.
func IsIdExistsValidation(db database.SNS) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		if sh, err := db.GetShorten(c.Context(), int64(id)); err != nil && sh.ID == 0 {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("id %d is not found", sh.ID),
			})
		}

		return c.Next()
	}
}
