package routes

import (
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/mdanialr/sns_backend/internal/service"
)

func SetupRoutes(app *fiber.App, conf *service.Config, fl io.Writer, db *database.Queries) {
	// Built-in fiber middleware
	app.Use(recover.New())
	app.Use(helmet.New())
	// Use log file only in production otherwise output it to stdout
	switch conf.EnvIsProd {
	case true:
		fConf := logger.Config{
			Format:     "[${time}] ${status} | ${method} - ${latency} - ${ip} | ${path}\n",
			TimeFormat: "02-Jan-2006 15:04:05",
			Output:     fl,
		}
		app.Use(logger.New(fConf))
	case false:
		app.Use(logger.New())
	}

	// This app's endpoints
	apiRoute := app.Group("/api")

	v1 := apiRoute.Group("/v1")

	sh := v1.Group("/shorten")
	ShortenRoutes(sh, db)

	sn := v1.Group("/send")
	SendRoutes(sn, conf, db)

	// Custom middleware AFTER endpoints
	//app.Use(api.DefaultRouteNotFound)
}
