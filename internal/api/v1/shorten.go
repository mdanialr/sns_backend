package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
)

// CreateShortenPayload json payload that need to be sent for sending request to this endpoint.
type CreateShortenPayload struct {
	database.CreateShortenParams
	IsReplace bool `json:"replace"`
}

// CreateShorten endpoint to create new shorten url.
func CreateShorten(db database.SNS) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var payload CreateShortenPayload
		if err := c.BodyParser(&payload); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to parse json payload: %s", err),
			})
		}

		// check whether provided Url already exist in database or not
		oldSh, _ := db.GetShortenByUrl(c.Context(), payload.Url)
		switch payload.IsReplace {
		case true:
			// if provided `replace` is true then delete it
			// if given Url really exist then it's Url field should not be empty and need to be deleted
			// if given Url does not exist then it's Url would be empty therefor does not need to be deleted
			if oldSh.Url == payload.Url {
				if err := db.DeleteShorten(c.Context(), oldSh.ID); err != nil {
					c.Status(fiber.StatusInternalServerError)
					return c.JSON(fiber.Map{
						"message": fmt.Sprintf("failed to delete old shorten data with url %s: %s", oldSh.Url, err),
					})
				}
			}
		case false:
			// if provided `replace` is false, then return error
			// data that founded from database should have non-zero ID
			if oldSh.ID != 0 {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"message": fmt.Sprintf("the provided url `%s` already exist", payload.Url),
				})
			}
		}

		sh, err := db.CreateShorten(c.Context(), payload.CreateShortenParams)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to create new shorten with the given payload: %s", err),
			})
		}

		c.Status(fiber.StatusCreated)
		return c.JSON(fiber.Map{
			"message": "successfully created",
			"data":    sh,
		})
	}
}
