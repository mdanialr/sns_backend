package send_handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mdanialr/sns_backend/internal/core/service/send_service/mocks"
	"github.com/spf13/viper"
)

type (
	sendRoutes struct {
		Index, Create, Update, Delete string
	}
	sendDeps struct {
		sendSvc *mocks.Mocksend_serviceIService
	}
	helperSetup struct {
		App *fiber.App
		Dep sendDeps
		R   sendRoutes
		V   *viper.Viper
	}
)

// setupJSONReq set up request instance and add JSON request header.
func (h *helperSetup) setupJSONReq(method, route string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, route, body)
	req.Header.Add("Content-Type", fiber.MIMEApplicationJSONCharsetUTF8)

	return req
}

func setupHelperTest(v *viper.Viper) *helperSetup {
	r := sendRoutes{
		Index:  "/send/",
		Create: "/send/create",
		Update: "/send/update",
		Delete: "/send/delete",
	}
	d := sendDeps{
		sendSvc: new(mocks.Mocksend_serviceIService),
	}

	return &helperSetup{
		App: fiber.New(),
		Dep: d,
		R:   r,
		V:   v,
	}
}

func defaultViper() *viper.Viper {
	v := viper.New()
	v.Set("jwt.secret", jwtSecret)
	return v
}

// createJWT return jwt token based on given duration and secret.
func createJWT(dur, secret string) string {
	return createJWTWithUser(dur, secret, secret)
}

// createJWTWithUser return jwt token based on given duration, secret also user in JWT Claims.
func createJWTWithUser(dur, secret, user string) string {
	d, _ := time.ParseDuration(dur)
	claims := jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(d).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(secret))
	return t
}
