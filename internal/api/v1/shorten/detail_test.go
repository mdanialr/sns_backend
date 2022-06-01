package api

import (
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

func TestGetShortenDetail(t *testing.T) {
	type resPassJSON struct {
		Msg database.Shorten `json:"message"`
	}

	t.Run("Should pass if the given id is exist in database and return the intended shorten data", func(t *testing.T) {
		sample := database.Shorten{
			ID:        11,
			Url:       "go",
			Target:    "https://go.dev",
			Permanent: true,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sns := mockdb.NewMockSNS(ctrl)
		sns.EXPECT().
			GetShorten(gomock.Any(), sample.ID).
			Times(1).
			Return(sample, nil)

		app := fiber.New()
		app.Get("/:id", GetShortenDetail(sns))

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", sample.ID), nil)
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		var r resPassJSON
		_ = json.NewDecoder(res.Body).Decode(&r)
		assert.Equal(t, r.Msg, sample)
	})

	testCases := []struct {
		name       string
		sample     int64
		buildStubs func(sns *mockdb.MockSNS, id int64)
		expectCode int
		expectMsg  any
	}{
		{
			name: "Should fail when database failed to get the intended shorten data, this error should has nothing to" +
				" do with our code but should be from database themself therefor return code 500",
			sample: -2,
			buildStubs: func(sns *mockdb.MockSNS, id int64) {
				sns.EXPECT().
					GetShorten(gomock.Any(), id).
					Times(1).
					Return(database.Shorten{}, fmt.Errorf("oops something goes wrong"))
			},
			expectCode: fiber.StatusInternalServerError,
			expectMsg:  "failed to retrieve a shorten data",
		},
		{
			name:   "Should fail if the given id is not exist in database",
			sample: 0,
			buildStubs: func(sns *mockdb.MockSNS, id int64) {
				sns.EXPECT().
					GetShorten(gomock.Any(), id).
					Times(1).
					Return(database.Shorten{}, sql.ErrNoRows)
			},
			expectCode: fiber.StatusNotFound,
			expectMsg:  "is not found in database",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sns := mockdb.NewMockSNS(ctrl)
			tc.buildStubs(sns, tc.sample)

			app := fiber.New()
			app.Get("/:id", GetShortenDetail(sns))

			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", tc.sample), nil)
			res, _ := app.Test(req)

			assert.Equal(t, tc.expectCode, res.StatusCode)
			var r JsonResponse
			_ = json.NewDecoder(res.Body).Decode(&r)
			assert.Contains(t, r.Msg, tc.expectMsg)
		})
	}
}
