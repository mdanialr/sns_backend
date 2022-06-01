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

func TestDeleteShorten(t *testing.T) {
	type resJSON struct {
		Msg string `json:"message"`
	}

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
			name: "Should fail when database failed to delete the intended shorten data, this error should has nothing to" +
				" do with our code but should be from database themself therefor return code 500",
			sample: 2,
			buildStubs: func(s *mockdb.MockSNS, id int64) {
				s.EXPECT().
					GetShorten(gomock.Any(), id).
					Times(1).
					Return(database.Shorten{ID: 12}, nil)

				s.EXPECT().
					DeleteShorten(gomock.Any(), id).
					Times(1).
					Return(sql.ErrNoRows)
			},
			expectCode: fiber.StatusInternalServerError,
			expectMsg:  "failed to delete a shorten data",
		},
		{
			name:   "Should pass when there are no errors in database operations which are get and delete a shorten data",
			sample: 12,
			buildStubs: func(s *mockdb.MockSNS, id int64) {
				s.EXPECT().
					GetShorten(gomock.Any(), id).
					Times(1).
					Return(database.Shorten{ID: 12}, nil)

				s.EXPECT().
					DeleteShorten(gomock.Any(), id).
					Times(1).
					Return(nil)
			},
			expectCode: fiber.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sns := mockdb.NewMockSNS(ctrl)
			tc.buildStubs(sns, tc.sample)

			app := fiber.New()
			app.Delete("/:id", DeleteShorten(sns))

			req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%d", tc.sample), nil)
			res, _ := app.Test(req)

			assert.Equal(t, tc.expectCode, res.StatusCode)
			var r resJSON
			_ = json.NewDecoder(res.Body).Decode(&r)
			fmt.Println(r.Msg)
			assert.Contains(t, r.Msg, tc.expectMsg)
		})
	}
}
