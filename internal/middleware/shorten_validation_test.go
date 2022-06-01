package middleware

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	mockdb "github.com/mdanialr/sns_backend/internal/database/mock"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortenValidation(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBuffer(tc.sampleBody))
			req.Header.Add("content-type", tc.sampleContentType)
			res, _ := app.Test(req)

			assert.Equal(t, tc.expectStatusCode, res.StatusCode)
			var r JsonResponse
			_ = json.NewDecoder(res.Body).Decode(&r)
			assert.Contains(t, r.Msg, tc.expectMsg)
		})
	}
}

func TestIsIdExistsValidation(t *testing.T) {
	testCases := []struct {
		name       string
		sample     int64
		buildStubs func(sns *mockdb.MockSNS, id int64)
		expectCode int
		expectMsg  string
	}{
		{
			name: "Should fail when database failed to get the intended shorten data, this error probably caused by" +
				" the ID is not exist in database therefor return 404 code",
			sample: 2,
			buildStubs: func(s *mockdb.MockSNS, id int64) {
				s.EXPECT().
					GetShorten(gomock.Any(), id).
					Times(1).
					Return(database.Shorten{ID: 0}, sql.ErrNoRows)
			},
			expectCode: fiber.StatusNotFound,
			expectMsg:  "is not found",
		},
		{
			name: "Should pass when using ID that indeed exists in the database and should pass this middleware and" +
				" get json response from the next handler",
			sample: 12,
			buildStubs: func(s *mockdb.MockSNS, id int64) {
				s.EXPECT().
					GetShorten(gomock.Any(), id).
					Times(1).
					Return(database.Shorten{ID: 12}, nil)
			},
			expectCode: fiber.StatusOK,
			expectMsg:  "success",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sns := mockdb.NewMockSNS(ctrl)
			tc.buildStubs(sns, tc.sample)

			app := fiber.New()
			app.Patch("/:id",
				IsIdExistsValidation(sns),
				func(c *fiber.Ctx) error {
					c.Status(fiber.StatusOK)
					return c.JSON(fiber.Map{"message": "success"})
				},
			)

			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/%d", tc.sample), nil)
			res, _ := app.Test(req)

			assert.Equal(t, tc.expectCode, res.StatusCode)
			var r JsonResponse
			_ = json.NewDecoder(res.Body).Decode(&r)
			assert.Contains(t, r.Msg, tc.expectMsg)
		})
	}
}
