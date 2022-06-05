package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/bytedance/sonic"
	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/mdanialr/sns_backend/internal/logger"
	"github.com/mdanialr/sns_backend/internal/routes"
	"github.com/mdanialr/sns_backend/internal/service"
	"github.com/spf13/viper"
)

func main() {
	// set default values first so it can be applied immediately
	viper.SetDefault("env", "dev")
	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", "7575")
	viper.SetDefault("log", "/tmp/")
	viper.SetDefault("upload", "/tmp/")

	appConfig, err := service.NewConfig("app", ".")
	if err != nil {
		log.Fatalln("failed to create new config:", err)
	}

	app, err := setup(appConfig)
	if err != nil {
		log.Fatalln("failed setup the app:", err)
	}

	// live reload on config file changes
	viper.OnConfigChange(func(_ fsnotify.Event) {
		if err := viper.Unmarshal(&appConfig); err != nil {
			logger.ErrL.Println("failed to reload and unmarshalling config:", err)
		}
	})
	viper.WatchConfig()

	// init custom app logger
	fl, err := os.OpenFile(appConfig.LogDir+"app-log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		log.Fatalln("failed to open|create fiber app log file:", err)
	}
	defer fl.Close()

	// setup database instance
	db, err := sql.Open(appConfig.DB.GetDriver(), appConfig.DB.GenerateConnectionString())
	if err != nil {
		log.Fatalln("failed to open connection to database:", err)
	}
	defer db.Close()
	dbInstance := database.New(db)

	routes.SetupRoutes(app, appConfig, fl, dbInstance)

	logger.InfL.Printf("listening on %s:%s\n", appConfig.Host, appConfig.PortNum)
	logger.ErrL.Fatalln(app.Listen(fmt.Sprintf("%s:%s", appConfig.Host, appConfig.PortNum)))
}

// setup prepare every necessary things before starting this app.
func setup(conf *service.Config) (*fiber.App, error) {
	conf.SanitizeEnv()
	conf.SanitizeDir()

	// init internal logging
	if err := logger.InitLogger(conf); err != nil {
		return nil, fmt.Errorf("failed to init internal logging: %v\n", err)
	}

	// if app in production, then use hostname from proxy instead
	var proxyHeader string
	if conf.EnvIsProd {
		proxyHeader = "X-Real-Ip"
	}

	app := fiber.New(fiber.Config{
		BodyLimit:             50 * 1024 * 1024,
		DisableStartupMessage: conf.EnvIsProd,
		ProxyHeader:           proxyHeader,
		JSONDecoder:           sonic.Unmarshal,
		JSONEncoder:           sonic.Marshal,
	})

	return app, nil
}
