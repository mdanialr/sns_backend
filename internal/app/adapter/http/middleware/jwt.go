package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtMiddleware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	cons "github.com/mdanialr/sns_backend/pkg/constant"
	resp "github.com/mdanialr/sns_backend/pkg/response"
	"github.com/spf13/viper"
)

// JWT middleware that use JSON Web Token as access token.
func JWT(v *viper.Viper) fiber.Handler {
	return jwtMiddleware.New(jwtMiddleware.Config{
		ContextKey:    "jwt",
		SigningMethod: "HS256",
		SigningKey:    []byte(v.GetString("jwt.secret")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return resp.ErrorCode(c, fiber.StatusUnauthorized, resp.WithErrMsg(cons.InvalidToken))
		},
		SuccessHandler: validateUser(v),
	})
}

// validateUser make sure user in JWT claim same as in when generating them
// which is using the JWT secret as the user.
func validateUser(v *viper.Viper) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jw := c.Locals("jwt")
		if jw == nil {
			return resp.ErrorCode(c, fiber.StatusInternalServerError, resp.WithErrMsg("empty user in local"))
		}
		tk := jw.(*jwt.Token)
		cl := tk.Claims.(jwt.MapClaims)
		if cl["user"] != v.GetString("jwt.secret") {
			return resp.Error(c, resp.WithErrMsg("unexpected user was found in jwt token"))
		}
		return c.Next()
	}
}
