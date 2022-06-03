package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	mockdb "github.com/mdanialr/sns_backend/internal/database/mock"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/stretchr/testify/assert"
)

func TestListShorten(t *testing.T) {
	type resPassJSON struct {
		Msg []database.Shorten `json:"message"`
	}
	sampleShortens := []database.Shorten{
		database.Shorten{
			ID:        1,
			Url:       "go",
			Target:    "https://go.dev",
			Permanent: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		database.Shorten{
			ID:        2,
			Url:       "dart",
			Target:    "https://pub.dev",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	t.Run("Should pass and returned list of shortens data should be as expected", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sns := mockdb.NewMockSNS(ctrl)
		sns.EXPECT().
			ListShorten(gomock.Any(), "url", database.ASC).
			Times(1).
			Return(sampleShortens, nil)

		app := fiber.New()
		app.Get("/shorten", ListShorten(sns))

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/shorten?sort=%s", "url"), nil)
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		var r resPassJSON
		_ = json.NewDecoder(res.Body).Decode(&r)
		for _, s := range r.Msg {
			for _, ss := range sampleShortens {
				if s.ID == ss.ID {
					assert.Equal(t, s.ID, ss.ID)
					assert.Equal(t, s.Url, ss.Url)
					assert.Equal(t, s.Target, ss.Target)
					assert.Equal(t, s.Permanent, ss.Permanent)
					assert.WithinDuration(t, s.CreatedAt, ss.CreatedAt, 0)
					assert.WithinDuration(t, s.UpdatedAt, ss.UpdatedAt, 0)
				}
			}
		}
	})

	testCases := []struct {
		name        string
		sampleQ     string
		sampleCol   string
		sampleOrder database.DBOrder
		buildStubs  func(sns *mockdb.MockSNS, c string, or database.DBOrder)
		expectCode  int
		expectMsg   string
	}{
		{
			name: "Should failed when database failed to listing all shorten data with 'url' param, this error should" +
				" has nothing to do with our code base but should be from database themself therefor return code 500",
			sampleQ:     "url",
			sampleCol:   "url",
			sampleOrder: database.ASC,
			buildStubs: func(s *mockdb.MockSNS, c string, or database.DBOrder) {
				s.EXPECT().
					ListShorten(gomock.Any(), c, or).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			expectCode: fiber.StatusInternalServerError,
			expectMsg:  "failed to listing all shorten data with 'url' sorting",
		},
		{
			name: "Should failed when database failed to listing all shorten data with 'date' param, this error should" +
				" has nothing to do with our code base but should be from database themself therefor return code 500",
			sampleQ:     "date",
			sampleCol:   "created_at",
			sampleOrder: database.DESC,
			buildStubs: func(s *mockdb.MockSNS, c string, or database.DBOrder) {
				s.EXPECT().
					ListShorten(gomock.Any(), c, or).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			expectCode: fiber.StatusInternalServerError,
			expectMsg:  "failed to listing all shorten data with 'date' sorting",
		},
		{
			name: "Should failed when database failed to listing all shorten data with 'latest' param, this error should" +
				" has nothing to do with our code base but should be from database themself therefor return code 500",
			sampleQ:     "latest",
			sampleCol:   "updated_at",
			sampleOrder: database.DESC,
			buildStubs: func(s *mockdb.MockSNS, c string, or database.DBOrder) {
				s.EXPECT().
					ListShorten(gomock.Any(), c, or).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			expectCode: fiber.StatusInternalServerError,
			expectMsg:  "failed to listing all shorten data with 'latest' sorting",
		},
		{
			name: "Should failed when database failed to listing all shorten data without any param, this error should" +
				" has nothing to do with our code base but should be from database themself therefor return code 500",
			sampleCol:   "id",
			sampleOrder: database.ASC,
			buildStubs: func(s *mockdb.MockSNS, c string, or database.DBOrder) {
				s.EXPECT().
					ListShorten(gomock.Any(), c, or).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			expectCode: fiber.StatusInternalServerError,
			expectMsg:  "failed to listing all shorten data without any sorting order",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sns := mockdb.NewMockSNS(ctrl)
			tc.buildStubs(sns, tc.sampleCol, tc.sampleOrder)

			app := fiber.New()
			app.Get("/shorten", ListShorten(sns))

			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/shorten?sort=%s", tc.sampleQ), nil)
			res, _ := app.Test(req)

			assert.Equal(t, tc.expectCode, res.StatusCode)
			var r JsonResponse
			_ = json.NewDecoder(res.Body).Decode(&r)
			assert.Contains(t, r.Msg, tc.expectMsg)
		})
	}
}
