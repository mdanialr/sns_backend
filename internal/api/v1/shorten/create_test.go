package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	mockdb "github.com/mdanialr/sns_backend/internal/database/mock"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShorten(t *testing.T) {
	shorten := database.Shorten{Url: "go", ID: 11}

	const testName = "Should fail when sending empty content-type and return code 400 then has expected error message"
	t.Run(testName, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sns := mockdb.NewMockSNS(ctrl)
		app := fiber.New()
		app.Post("/", CreateShorten(sns))

		req := httptest.NewRequest(fiber.MethodPost, "/", nil)
		req.Header.Add("content-type", "")
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		var r JsonResponse
		_ = json.NewDecoder(res.Body).Decode(&r)
		assert.Contains(t, r.Msg, "failed to parse json payload")
	})

	testCases := []struct {
		name       string
		sampleBody CreateShortenPayload
		buildStubs func(sns *mockdb.MockSNS)
		expectCode int
		expectMsg  string
	}{
		{
			name: "Should fail when database failed to delete already exist data, this error should has nothing to do" +
				" with our code but should be from database themself therefor return code 500",
			sampleBody: CreateShortenPayload{
				IsReplace: true,
				CreateShortenParams: database.CreateShortenParams{
					Url:       "go",
					Target:    "https://go.dev",
					Permanent: true,
				},
			},
			buildStubs: func(sns *mockdb.MockSNS) {
				sns.EXPECT().
					GetShortenByUrl(gomock.Any(), gomock.Eq(shorten.Url)).
					Times(1).
					Return(shorten, nil)

				sns.EXPECT().
					DeleteShorten(gomock.Any(), gomock.Eq(shorten.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			expectCode: fiber.StatusInternalServerError,
			expectMsg:  "failed to delete",
		},
		{
			name: "Should fail when provided Url already exist in database but provided `replace` field is false",
			sampleBody: CreateShortenPayload{
				IsReplace: false,
				CreateShortenParams: database.CreateShortenParams{
					Url:       "go",
					Target:    "https://go.dev",
					Permanent: true,
				},
			},
			buildStubs: func(sns *mockdb.MockSNS) {
				sns.EXPECT().
					GetShortenByUrl(gomock.Any(), gomock.Eq(shorten.Url)).
					Times(1).
					Return(shorten, nil)
			},
			expectCode: fiber.StatusBadRequest,
			expectMsg:  "already exist",
		},
		{
			name: "Should fail when database failed to create or save provided data, this error should has nothing to do" +
				" with our code but should be from database themself therefor return code 500",
			sampleBody: CreateShortenPayload{
				IsReplace: false,
				CreateShortenParams: database.CreateShortenParams{
					Url:       "go",
					Target:    "https://go.dev",
					Permanent: true,
				},
			},
			buildStubs: func(sns *mockdb.MockSNS) {
				sns.EXPECT().
					GetShortenByUrl(gomock.Any(), gomock.Eq(shorten.Url)).
					Times(1).
					Return(database.Shorten{}, nil)

				sns.EXPECT().
					CreateShorten(gomock.Any(), database.CreateShortenParams{
						Url:       "go",
						Target:    "https://go.dev",
						Permanent: true,
					}).
					Times(1).
					Return(database.Shorten{}, sql.ErrNoRows)
			},
			expectCode: fiber.StatusInternalServerError,
			expectMsg:  "failed to create",
		},
		{
			name: "Should pass when all database operations not returning any error",
			sampleBody: CreateShortenPayload{
				IsReplace: false,
				CreateShortenParams: database.CreateShortenParams{
					Url:       "go",
					Target:    "https://go.dev",
					Permanent: true,
				},
			},
			buildStubs: func(sns *mockdb.MockSNS) {
				sns.EXPECT().
					GetShortenByUrl(gomock.Any(), gomock.Eq(shorten.Url)).
					Times(1).
					Return(database.Shorten{}, nil)

				sns.EXPECT().
					CreateShorten(gomock.Any(), database.CreateShortenParams{
						Url:       "go",
						Target:    "https://go.dev",
						Permanent: true,
					}).
					Times(1).
					Return(database.Shorten{
						ID:        14,
						Url:       "go",
						Target:    "https://go.dev",
						Permanent: true,
					}, nil)
			},
			expectCode: fiber.StatusCreated,
			expectMsg:  "successfully created",
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sns := mockdb.NewMockSNS(ctrl)
			tc.buildStubs(sns)

			app := fiber.New()
			app.Post("/", CreateShorten(sns))

			b, err := json.Marshal(tc.sampleBody)
			require.NoErrorf(t, err, "failed to marshall json body in #%d test:", i+1)
			req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBuffer(b))
			req.Header.Add("content-type", fiber.MIMEApplicationJSON)
			res, _ := app.Test(req)

			assert.Equal(t, tc.expectCode, res.StatusCode)
			var r JsonResponse
			_ = json.NewDecoder(res.Body).Decode(&r)
			assert.Contains(t, r.Msg, tc.expectMsg)
		})
	}
}
