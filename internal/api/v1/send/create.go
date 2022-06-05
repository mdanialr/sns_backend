package api

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/mdanialr/sns_backend/internal/service"
)

// CreateSendPayload json payload that need to be sent for sending request to this endpoint.
type CreateSendPayload struct {
	Url       string `form:"url"`
	Permanent bool   `form:"permanent"`
	IsReplace bool   `form:"replace"`
}

// CreateSend endpoint to create new Send file.
func CreateSend(conf *service.Config, db database.SNS) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var payload CreateSendPayload
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

		f, err := c.FormFile("file")
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to get file instance from multipart request: %s", err),
			})
		}

		// rename file to make sure it's unique in local disk
		fileName := strings.TrimSuffix(f.Filename, filepath.Ext(f.Filename))
		newName := fmt.Sprintf("%s_%s%s",
			fileName,
			time.Now().Format("06-Jan-02_15-04-05.999"),
			filepath.Ext(f.Filename),
		)
		// save file to local
		if err := c.SaveFile(f, fmt.Sprintf("%s%s", conf.UploadDir, newName)); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to save uploaded file to local disk: %s", err),
			})
		}

		data := database.CreateSendParams{
			Url:       payload.Url,
			File:      newName,
			Size:      service.FormatBytesToHumanString(f.Size),
			Permanent: payload.Permanent,
		}

		sn, err := db.CreateSend(c.Context(), data)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to create new send with the given payload: %s", err),
			})
		}

		c.Status(fiber.StatusCreated)
		return c.JSON(fiber.Map{
			"message": "successfully created",
			"data":    sn,
		})
	}
}
