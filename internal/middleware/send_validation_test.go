package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSendValidation(t *testing.T) {
	app := fiber.New()
	app.Post("/",
		CreateSendValidation,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "success"})
		},
	)

	const name = "Should fail when sending empty or content-type other than multipart and return 400 code"
	t.Run(name, func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, "/", nil)
		req.Header.Add("content-type", fiber.MIMEApplicationJSON)
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		var r JsonResponse
		_ = json.NewDecoder(res.Body).Decode(&r)
		assert.Contains(t, r.Msg, "failed to parse multipart payload")
	})

	testCases := []struct {
		name       string
		sample     string
		expectCode int
		expectMsg  string
	}{
		{
			name:       "Should fail when not sending required `file` field in multipart request",
			sample:     "upload",
			expectCode: fiber.StatusBadRequest,
			expectMsg:  "`file` field is required",
		},
		{
			name:       "Should pass when sending the required `file` field in multipart request",
			sample:     "file",
			expectCode: fiber.StatusOK,
			expectMsg:  "success",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := new(bytes.Buffer)
			wr := multipart.NewWriter(body)
			fl, err := os.Open(fakeFilePath)
			require.NoError(t, err, "failed opening fake file path")
			defer fl.Close()

			form, err := wr.CreateFormFile(tc.sample, path.Base(fakeFilePath))
			require.NoError(t, err)

			if _, err := io.Copy(form, fl); err != nil {
				require.NoError(t, err, "failed copy fake file to multipart")
			}

			require.NoError(t, wr.Close(), "failed to close multipart writer")

			req := httptest.NewRequest(fiber.MethodPost, "/", body)
			req.Header.Add("content-type", wr.FormDataContentType())
			res, _ := app.Test(req)

			assert.Equal(t, tc.expectCode, res.StatusCode)
			var r JsonResponse
			_ = json.NewDecoder(res.Body).Decode(&r)
			assert.Contains(t, r.Msg, tc.expectMsg)
		})
	}

	t.Cleanup(func() {
		os.Remove(fakeFilePath)
	})
}
