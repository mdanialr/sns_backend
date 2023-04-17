package auth_handler_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/mdanialr/sns_backend/internal/app/adapter/http/auth_handler"
	"github.com/mdanialr/sns_backend/internal/core/service/otp_service/mocks"
	"github.com/mdanialr/sns_backend/internal/requests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Login(t *testing.T) {
	const sampleJWToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODAzNjQyOTAsInVzZXIiOiJSUm1QVEdpTjVsa0pSUFVjb2o4dnZIM2xYU0Ztbm5oUWZ5UjRLazdhcUJGM3hQVHFoYSJ9.LqLSf9AdXErRki3n-JuelmEl8PXB98PSzlHChayExMQ"

	testCases := []struct {
		name           string
		payload        io.Reader
		setup          func(*mocks.Mockotp_serviceIService)
		expectCode     int
		expectResponse string
	}{
		{
			name: "Given wrong key that should be 'code' in request payload and string otp asdasd should return error " +
				" message Invalid Payload and validation error required",
			payload:        bytes.NewBufferString(`{"otp":"asdasd"}`),
			setup:          func(_ *mocks.Mockotp_serviceIService) {},
			expectCode:     http.StatusBadRequest,
			expectResponse: `{"status":"FAILED","message":"Invalid Payload","detail":[{"name":"code","message":"required"}]}`,
		},
		{
			name: "Given invalid request and string otp asdasd should return error message Invalid Payload and validation" +
				" error should be numeric",
			payload:        bytes.NewBufferString(`{"code":"asdasd"}`),
			setup:          func(_ *mocks.Mockotp_serviceIService) {},
			expectCode:     http.StatusBadRequest,
			expectResponse: `{"status":"FAILED","message":"Invalid Payload","detail":[{"name":"code","message":"should be numeric"}]}`,
		},
		{
			name: "Given valid request otp 123 and the length is three chars should return error message Invalid " +
				"Payload and validation error length should be 6",
			payload:        bytes.NewBufferString(`{"code":"123"}`),
			setup:          func(_ *mocks.Mockotp_serviceIService) {},
			expectCode:     http.StatusBadRequest,
			expectResponse: `{"status":"FAILED","message":"Invalid Payload","detail":[{"name":"code","message":"length should be 6"}]}`,
		},
		{
			name: "Given valid request otp 1234567 and the length is seven chars should return error message Invalid " +
				"Payload and validation error length should be 6",
			payload:        bytes.NewBufferString(`{"code":"1234567"}`),
			setup:          func(_ *mocks.Mockotp_serviceIService) {},
			expectCode:     http.StatusBadRequest,
			expectResponse: `{"status":"FAILED","message":"Invalid Payload","detail":[{"name":"code","message":"length should be 6"}]}`,
		},
		{
			name:    "Given invalid otp 123456 should return error message Invalid OTP",
			payload: bytes.NewBufferString(`{"code":"123456"}`),
			setup: func(svc *mocks.Mockotp_serviceIService) {
				svc.EXPECT().
					ValidateOTP(mock.Anything, &requests.OTP{Code: "123456"}).
					Return(false).
					Once()
			},
			expectCode:     http.StatusBadRequest,
			expectResponse: `{"status":"FAILED","message":"Invalid OTP"}`,
		},
		{
			name: "Given valid otp 654321 but failed to create jwt token should return error message from service" +
				"layer dependency and status code Internal Server Error",
			payload: bytes.NewBufferString(`{"code":"654321"}`),
			setup: func(svc *mocks.Mockotp_serviceIService) {
				svc.EXPECT().
					ValidateOTP(mock.Anything, &requests.OTP{Code: "654321"}).
					Return(true).
					Once()
				svc.EXPECT().
					GetJWT(mock.Anything).
					Return("", errors.New("failed to sign token")).
					Once()
			},
			expectCode:     http.StatusInternalServerError,
			expectResponse: `{"status":"FAILED","message":"failed to sign token"}`,
		},
		{
			name: "Given valid otp 654321 and successfully creating the JWT token should return that token and " +
				"status code OK",
			payload: bytes.NewBufferString(`{"code":"654321"}`),
			setup: func(svc *mocks.Mockotp_serviceIService) {
				svc.EXPECT().
					ValidateOTP(mock.Anything, &requests.OTP{Code: "654321"}).
					Return(true).
					Once()
				svc.EXPECT().
					GetJWT(mock.Anything).
					Return(sampleJWToken, nil).
					Once()
			},
			expectCode:     http.StatusOK,
			expectResponse: `{"status":"SUCCESS","data":{"token":"` + sampleJWToken + `"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// run necessary setup
			h := setupHelperTest()
			auth_handler.New(h.App, h.Dep.otpSvc)
			tc.setup(h.Dep.otpSvc)

			// setup request
			req := h.setupJSONReq(http.MethodPost, h.R.Otp, tc.payload)
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
