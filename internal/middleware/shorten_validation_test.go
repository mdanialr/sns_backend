package middleware

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortenValidation(t *testing.T) {
	type resJSON struct {
		Msg string `json:"message"`
	}
	app := fiber.New()
	app.Post("/",
		CreateShortenValidation,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "success"})
		},
	)

	testCases := []struct {
		name              string
		sampleBody        []byte
		sampleContentType string
		expectStatusCode  int
		expectMsg         string
	}{
		{
			name:             "Should fail when sending empty content-type and return code 400",
			expectStatusCode: fiber.StatusBadRequest,
			expectMsg:        "failed to parse json payload",
		},
		{
			name: "Should fail and return code 400 when sending empty request body even with the right" +
				" content-type",
			sampleContentType: fiber.MIMEApplicationJSON,
			expectStatusCode:  fiber.StatusBadRequest,
			expectMsg:         "failed to parse json payload",
		},
		{
			name:              "Should fail when sending required `target` field with empty value and return code 400",
			sampleBody:        []byte(`{"target": ""}`),
			sampleContentType: fiber.MIMEApplicationJSON,
			expectStatusCode:  fiber.StatusBadRequest,
			expectMsg:         "`target` field should not be empty",
		},
		{
			name: "Should pass and return 200 from next handler when sending valid json request and has" +
				" non empty value in `target` field",
			sampleBody:        []byte(`{"target": "https://go.dev"}`),
			sampleContentType: fiber.MIMEApplicationJSON,
			expectStatusCode:  fiber.StatusOK,
			expectMsg:         "success",
		},
		{
			name: "Should pass and return 200 from next handler and `url` field should not be empty after" +
				" passing this middleware",
			sampleBody:        []byte(`{"target": "https://go.dev", "url": ""}`),
			sampleContentType: fiber.MIMEApplicationJSON,
			expectStatusCode:  fiber.StatusOK,
			expectMsg:         "success",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBuffer(tc.sampleBody))
			req.Header.Add("content-type", tc.sampleContentType)
			res, _ := app.Test(req)

			assert.Equal(t, tc.expectStatusCode, res.StatusCode)
			var r resJSON
			_ = json.NewDecoder(res.Body).Decode(&r)
			assert.Contains(t, r.Msg, tc.expectMsg)
		})
	}
}
