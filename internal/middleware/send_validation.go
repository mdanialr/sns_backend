package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	api "github.com/mdanialr/sns_backend/internal/api/v1/send"
)

// CreateSendValidation middleware that handle necessary validation for creating new Send file.
func CreateSendValidation(c *fiber.Ctx) error {
	var payload api.CreateSendPayload
	if err := c.BodyParser(&payload); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("failed to parse multipart payload: %s", err),
		})
	}

	// 1st validation: `file` field is required
	if _, err := c.FormFile("file"); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "`file` field is required",
		})
	}

	return c.Next()
}
