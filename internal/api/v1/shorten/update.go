package api

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
)

// UpdateShortenPayload json payload that need to be sent for sending request to this endpoint.
type UpdateShortenPayload struct {
	CreateShortenPayload
}

// UpdateShorten update a shorten data using given ID in the url.
func UpdateShorten(db database.SNS) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		var payload UpdateShortenPayload
		if err := c.BodyParser(&payload); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to parse json payload: %s", err),
			})
		}

		var data database.UpdateShortenParams
		// update the updated_at field
		data.UpdatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		data.ID = int64(id)
		data.Url = payload.Url
		data.Target = payload.Target
		data.Permanent = payload.Permanent

		newData, err := db.UpdateShorten(c.Context(), data)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to update a shorten data with id %d and error: %s", id, err),
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": "successfully updated",
			"data":    newData,
		})
	}
}
