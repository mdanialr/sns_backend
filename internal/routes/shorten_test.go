package routes

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/stretchr/testify/require"
)

func TestShortenRoutes(t *testing.T) {
	const appLog = "/tmp/app-log"
	app := fiber.New()
	db := database.Queries{}
	fl, err := os.OpenFile(appLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	require.NoError(t, err)
	defer fl.Close()

	ShortenRoutes(app, &db)

	t.Cleanup(func() {
		os.Remove(appLog)
	})
}
