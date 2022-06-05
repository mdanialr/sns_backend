package routes

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/mdanialr/sns_backend/internal/service"
	"github.com/stretchr/testify/require"
)

func TestSendRoutes(t *testing.T) {
	const appLog = "/tmp/app-log"
	app := fiber.New()
	db := database.Queries{}
	conf := service.Config{}
	fl, err := os.OpenFile(appLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	require.NoError(t, err)
	defer fl.Close()

	SendRoutes(app, &conf, &db)

	t.Cleanup(func() {
		os.Remove(appLog)
	})
}
