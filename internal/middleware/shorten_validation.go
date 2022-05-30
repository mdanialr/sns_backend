package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/sns_backend/internal/api/v1/shorten"
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
