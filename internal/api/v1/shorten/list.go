package api

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
)

// ListShorten listing all shorten data, could be sorted by column name and the order could be ASC or DESC.
func ListShorten(db database.SNS) func(*fiber.Ctx) error {
	var (
		shortens []database.Shorten
		err      error
	)

	return func(c *fiber.Ctx) error {
		sort := c.Query("sort")

		switch sort {
		case "url":
			shortens, err = db.ListShorten(c.Context(), "url", database.ASC)
			if err != nil {
				err = errors.New(fmt.Sprintf("failed to listing all shorten data with 'url' sorting: %s", err))
			}
		case "date":
			shortens, err = db.ListShorten(c.Context(), "created_at", database.DESC)
			if err != nil {
				err = errors.New(fmt.Sprintf("failed to listing all shorten data with 'date' sorting: %s", err))
			}
		case "latest":
			shortens, err = db.ListShorten(c.Context(), "updated_at", database.DESC)
			if err != nil {
				err = errors.New(fmt.Sprintf("failed to listing all shorten data with 'latest' sorting: %s", err))
			}
		default:
			shortens, err = db.ListShorten(c.Context(), "id", database.ASC)
			if err != nil {
				err = errors.New(fmt.Sprintf("failed to listing all shorten data without any sorting order: %s", err))
			}
		}

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"message": err.Error()})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{"message": shortens})
	}
}
