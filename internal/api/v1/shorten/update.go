package api

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/mdanialr/sns_backend/internal/service"
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
		// validation: if `url` field is not provided then substitute it with random string
		if payload.Url == "" {
			payload.Url = service.RandomString(8)
		}

		// check whether provided Url for update already exist in database or not
		oldSh, _ := db.GetShortenByUrl(c.Context(), payload.Url)
		switch payload.IsReplace {
		case true:
			// if provided `replace` is true then delete it.
			// if founded data's ID with the target of update's ID is the same then do not need to delete it.
			// otherwise if the ID of both instance is not same then delete the other data.
			if oldSh.ID != int64(id) {
				if err := db.DeleteShorten(c.Context(), oldSh.ID); err != nil {
					c.Status(fiber.StatusInternalServerError)
					return c.JSON(fiber.Map{
						"message": fmt.Sprintf("failed to delete old shorten data with url %s: %s", oldSh.Url, err),
					})
				}
			}
		case false:
			// if provided `replace` is false, then return error that the provided Url already exist.
			if oldSh.ID != 0 {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"message": fmt.Sprintf("the provided url `%s` already exist", payload.Url),
				})
			}
		}

		var data database.UpdateShortenParams
		// update the updated_at field
		data.UpdatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		data.ID = oldSh.ID
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
