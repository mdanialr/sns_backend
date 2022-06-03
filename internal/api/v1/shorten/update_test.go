package api

import (
	"bytes"
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
	"github.com/stretchr/testify/require"
)

func TestUpdateShorten(t *testing.T) {
	shortenInDatabase := database.Shorten{
		ID:        11,
		Url:       "go",
		Target:    "https://pub.dev",
		Permanent: false,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	newShorten := database.UpdateShortenParams{
		ID:        11,
		Url:       shortenInDatabase.Url,
		Target:    "https://go.dev",
		Permanent: false,
	}

	const testName = "Should fail when sending empty content-type and return code 400 and has expected error message"
	t.Run(testName, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sns := mockdb.NewMockSNS(ctrl)
		app := fiber.New()
		app.Patch("/:id", UpdateShorten(sns))

		req := httptest.NewRequest(fiber.MethodPatch, "/23", nil)
		req.Header.Add("content-type", "")
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		var r JsonResponse
		_ = json.NewDecoder(res.Body).Decode(&r)
		assert.Contains(t, r.Msg, "failed to parse json payload")
	})

	testCases := []struct {
		name       string
		sample     int64
		sampleBody UpdateShortenPayload
		buildStubs func(sns *mockdb.MockSNS)
		expectCode int
		expectMsg  string
	}{
		{
			name: "Should fail when database failed to delete already exist data, this error should has nothing to do" +
				" with our code but should be from database themself therefor return code 500",
			sample: 14,
			sampleBody: UpdateShortenPayload{
				CreateShortenPayload: CreateShortenPayload{
					IsReplace: true,
					CreateShortenParams: database.CreateShortenParams{
						Url:       "go",
						Target:    "https://go.dev",
						Permanent: true,
					},
				},
			},
			buildStubs: func(sns *mockdb.MockSNS) {
				sns.EXPECT().
					GetShortenByUrl(gomock.Any(), gomock.Eq(shortenInDatabase.Url)).
					Times(1).
					Return(shortenInDatabase, nil)

				sns.EXPECT().
					DeleteShorten(gomock.Any(), gomock.Eq(shortenInDatabase.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			expectCode: fiber.StatusInternalServerError,
			expectMsg:  "failed to delete",
		},
		{
			name:   "Should fail when provided Url already exist in database but provided `replace` field is false",
			sample: 15,
			sampleBody: UpdateShortenPayload{
				CreateShortenPayload: CreateShortenPayload{
					IsReplace: false,
					CreateShortenParams: database.CreateShortenParams{
						Url:       "go",
						Target:    "https://go.dev",
						Permanent: true,
					},
				},
			},
			buildStubs: func(sns *mockdb.MockSNS) {
				sns.EXPECT().
					GetShortenByUrl(gomock.Any(), gomock.Eq(shortenInDatabase.Url)).
					Times(1).
					Return(shortenInDatabase, nil)
			},
			expectCode: fiber.StatusBadRequest,
			expectMsg:  "already exist",
		},
		{
			name: "Should fail when database failed to update the intended Shorten data, this error should has" +
				" nothing to do with our code but should be from database themself therefor return code 500",
			sample: shortenInDatabase.ID,
			sampleBody: UpdateShortenPayload{
				CreateShortenPayload: CreateShortenPayload{
					IsReplace: true,
					CreateShortenParams: database.CreateShortenParams{
						Url:       "go",
						Target:    "https://go.dev",
						Permanent: false,
					},
				},
			},
			buildStubs: func(sns *mockdb.MockSNS) {
				sns.EXPECT().
					GetShortenByUrl(gomock.Any(), gomock.Eq(shortenInDatabase.Url)).
					Times(1).
					Return(shortenInDatabase, nil)

				sns.EXPECT().
					UpdateShorten(gomock.Any(), gomock.AssignableToTypeOf(newShorten)).
					Times(1).
					Return(shortenInDatabase, sql.ErrNoRows)

			},
			expectCode: fiber.StatusInternalServerError,
			expectMsg:  "failed to update",
		},
		{
			name:   "Should pass when all database operations not returning any error",
			sample: shortenInDatabase.ID,
			sampleBody: UpdateShortenPayload{
				CreateShortenPayload: CreateShortenPayload{
					IsReplace: true,
					CreateShortenParams: database.CreateShortenParams{
						Url:       "go",
						Target:    "https://go.dev",
						Permanent: false,
					},
				},
			},
			buildStubs: func(sns *mockdb.MockSNS) {
				sns.EXPECT().
					GetShortenByUrl(gomock.Any(), gomock.Eq(shortenInDatabase.Url)).
					Times(1).
					Return(shortenInDatabase, nil)

				sns.EXPECT().
					UpdateShorten(gomock.Any(), gomock.AssignableToTypeOf(newShorten)).
					Times(1).
					Return(shortenInDatabase, nil)

			},
			expectCode: fiber.StatusOK,
			expectMsg:  "successfully updated",
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sns := mockdb.NewMockSNS(ctrl)
			tc.buildStubs(sns)

			app := fiber.New()
			app.Patch("/:id", UpdateShorten(sns))

			b, err := json.Marshal(tc.sampleBody)
			require.NoErrorf(t, err, "failed to marshall json body in #%d test:", i+1)
			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/%d", tc.sample), bytes.NewBuffer(b))
			req.Header.Add("content-type", fiber.MIMEApplicationJSON)
			res, _ := app.Test(req)

			assert.Equal(t, tc.expectCode, res.StatusCode)
			var r JsonResponse
			_ = json.NewDecoder(res.Body).Decode(&r)
			assert.Contains(t, r.Msg, tc.expectMsg)
		})
	}
}
