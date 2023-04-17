package auth_handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/sns_backend/internal/core/service/otp_service/mocks"
)

type (
	authRoutes struct {
		Otp string
	}
	authDeps struct {
		otpSvc *mocks.Mockotp_serviceIService
	}
	helperSetup struct {
		App *fiber.App
		Dep authDeps
		R   authRoutes
	}
)

// setupJSONReq set up request instance and add JSON request header.
func (h *helperSetup) setupJSONReq(method, route string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, route, body)
	req.Header.Add("Content-Type", fiber.MIMEApplicationJSONCharsetUTF8)

	return req
}

func setupHelperTest() *helperSetup {
	r := authRoutes{
		Otp: "/auth/otp",
	}
	d := authDeps{
		otpSvc: new(mocks.Mockotp_serviceIService),
	}

	return &helperSetup{
		App: fiber.New(),
		Dep: d,
		R:   r,
	}
}
