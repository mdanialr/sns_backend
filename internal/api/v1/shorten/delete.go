package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
)

// DeleteShorten delete a shorten data using given ID in the url.
func DeleteShorten(db database.SNS) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		if sh, err := db.GetShorten(c.Context(), int64(id)); err != nil && sh.ID == 0 {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("id %d is not found", sh.ID),
			})
		}

		if err := db.DeleteShorten(c.Context(), int64(id)); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to delete a shorten data with id %d and error: %s", id, err),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
