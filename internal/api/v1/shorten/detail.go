package api

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
)

// GetShortenDetail endpoint to get the detail of a row from Shorten data.
func GetShortenDetail(db database.SNS) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		sh, err := db.GetShorten(c.Context(), int64(id))
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				c.Status(fiber.StatusNotFound)
				return c.JSON(fiber.Map{
					"message": fmt.Sprintf("shorten data with id: %d is not found in database", id),
				})
			default:
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": fmt.Sprintf("failed to retrieve a shorten data with id %d: %s", id, err),
				})
			}
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": sh,
		})
	}
}
