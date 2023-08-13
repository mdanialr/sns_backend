package send_handler_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/mdanialr/sns_backend/internal/app/adapter/http/send_handler"
	"github.com/mdanialr/sns_backend/internal/core/service/send_service/mocks"
	"github.com/mdanialr/sns_backend/internal/responses"
	"github.com/mdanialr/sns_backend/pkg/helper"
	paginate "github.com/mdanialr/sns_backend/pkg/pagination"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	jwtSecret       = "secret"
	jwtDur          = "1m" // 1 minute is enough for every test run
	jwtTokenExpired = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODE2OTUxNTUsInVzZXIiOiJzZWNyZXQifQ.EF_AL-_p8rAF5S2i1NrgjP1Wq7Fovz4wREa-Aj5djGQ"
)

func TestSendHandler_Index(t *testing.T) {
	const sampleOk = `{"status":"SUCCESS","data":[{"id":82,"url":"zoom-mixer","description":"","send":"upload/VWTIrBGnOCEwxuXLzjYcogl6eitF4f.png","file_size":"36.7KB","permanent":true,"created_at":"2021-04-28T18:27:45+07:00","updated_at":"2021-04-28T18:27:45+07:00"}],"meta":{"per_page":1,"current_page":1,"next_page":2,"total_page":20}}`

	testCases := []struct {
		name           string
		setupV         func() *viper.Viper
		jwtToken       string
		setup          func(*mocks.Mocksend_serviceIService)
		payload        io.Reader
		expectCode     int
		expectResponse string
	}{
		{
			name: "Given empty authorization header should return error message Invalid or Expired Token and" +
				"status code Unauthorized",
			setupV:         defaultViper,
			jwtToken:       "",
			setup:          func(_ *mocks.Mocksend_serviceIService) {},
			expectCode:     http.StatusUnauthorized,
			expectResponse: `{"status":"FAILED","message":"Invalid or Expired Token"}`,
		},
		{
			name: "Given invalid random jwt token in authorization header should return error message Invalid " +
				"or Expired Token and status code Unauthorized",
			setupV:         defaultViper,
			jwtToken:       "random",
			setup:          func(_ *mocks.Mocksend_serviceIService) {},
			expectCode:     http.StatusUnauthorized,
			expectResponse: `{"status":"FAILED","message":"Invalid or Expired Token"}`,
		},
		{
			name: "Given valid jwt token but has different signing secret in authorization header should return " +
				"error message Invalid or Expired Token and status code Unauthorized",
			setupV:         defaultViper,
			jwtToken:       createJWT(jwtDur, "not"+jwtSecret),
			setup:          func(_ *mocks.Mocksend_serviceIService) {},
			expectCode:     http.StatusUnauthorized,
			expectResponse: `{"status":"FAILED","message":"Invalid or Expired Token"}`,
		},
		{
			name: "Given valid jwt token with right signing secret but already expired in authorization header " +
				"should return error message Invalid or Expired Token and status code Unauthorized",
			setupV:         defaultViper,
			jwtToken:       jwtTokenExpired,
			setup:          func(_ *mocks.Mocksend_serviceIService) {},
			expectCode:     http.StatusUnauthorized,
			expectResponse: `{"status":"FAILED","message":"Invalid or Expired Token"}`,
		},
		{
			name: "Given valid jwt token with right signing secret but has different user inside the JWT Claims " +
				"in authorization header should return error message unexpected user was found in jwt token and " +
				"status code Bad Request",
			setupV:         defaultViper,
			jwtToken:       createJWTWithUser(jwtDur, jwtSecret, "user"),
			setup:          func(_ *mocks.Mocksend_serviceIService) {},
			expectCode:     http.StatusBadRequest,
			expectResponse: `{"status":"FAILED","message":"unexpected user was found in jwt token"}`,
		},
		{
			name: "Given right jwt token but failed to retrieve data from dependency should return error message " +
				"from service layer dependency and status code Bad Request",
			setupV:   defaultViper,
			jwtToken: createJWT(jwtDur, jwtSecret),
			setup: func(svc *mocks.Mocksend_serviceIService) {
				svc.EXPECT().
					Index(mock.Anything, mock.Anything).
					Return(nil, errors.New("failed to get send data")).
					Once()
			},
			expectCode:     http.StatusBadRequest,
			expectResponse: `{"status":"FAILED","message":"failed to get send data"}`,
		},
		{
			name: "Given valid jwt token and successfully to retrieve data from dependency should return the data " +
				"and status code OK",
			setupV:   defaultViper,
			jwtToken: createJWT(jwtDur, jwtSecret),
			setup: func(svc *mocks.Mocksend_serviceIService) {
				cr, _ := time.Parse(time.RFC3339, "2021-04-28T18:27:45+07:00")
				up, _ := time.Parse(time.RFC3339, "2021-04-28T18:27:45+07:00")
				obj := &responses.SendIndexResponse{
					Data: []*responses.SendResponse{{
						ID:          82,
						Url:         "zoom-mixer",
						Send:        helper.Ptr("upload/VWTIrBGnOCEwxuXLzjYcogl6eitF4f.png"),
						FileSize:    helper.Ptr("36.7KB"),
						IsPermanent: helper.Ptr(true),
						CreatedAt:   &cr,
						UpdatedAt:   &up,
					}},
					Pagination: &paginate.M{
						Limit:     1,
						Page:      1,
						Next:      2,
						TotalPage: 20,
					},
				}
				svc.EXPECT().
					Index(mock.Anything, mock.Anything).
					Return(obj, nil).
					Once()
			},
			expectCode:     http.StatusOK,
			expectResponse: sampleOk,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// run necessary setup
			h := setupHelperTest(tc.setupV())
			send_handler.New(h.App, h.V, h.Dep.sendSvc)
			tc.setup(h.Dep.sendSvc)

			// setup request payload
			req := h.setupJSONReq(http.MethodGet, h.R.Index, tc.payload)

			req.Header.Add("Authorization", "Bearer "+tc.jwtToken)
			res, _ := h.App.Test(req)
			defer res.Body.Close()

			assert.Equal(t, tc.expectCode, res.StatusCode)

			// assert the response payload
			var resp bytes.Buffer
			resp.ReadFrom(res.Body)
			assert.Equal(t, tc.expectResponse, resp.String())
		})
	}
}
