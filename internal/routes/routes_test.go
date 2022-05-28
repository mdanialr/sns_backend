package routes

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/mdanialr/sns_backend/internal/service"
	"github.com/stretchr/testify/require"
)

func TestSetupRoutes(t *testing.T) {
	const appLog = "/tmp/app-log"
	app := fiber.New()
	db := database.Queries{}
	fl, err := os.OpenFile(appLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	require.NoError(t, err)
	defer fl.Close()

	t.Run("1# Success test even there is no value supplied because there is no required value", func(t *testing.T) {
		conf := service.Config{}
		SetupRoutes(app, &conf, fl, &db)
	})

	t.Run("2# Success test with only one or more supplied value", func(t *testing.T) {
		conf := service.Config{EnvIsProd: true}
		SetupRoutes(app, &conf, fl, &db)
	})

	t.Cleanup(func() {
		os.Remove(appLog)
	})
}
