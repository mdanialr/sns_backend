package server

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	"github.com/mdanialr/sns_backend/internal/app"
	conf "github.com/mdanialr/sns_backend/pkg/config"
	gormLogger "github.com/mdanialr/sns_backend/pkg/gorm"
	"github.com/mdanialr/sns_backend/pkg/logger"
	"github.com/mdanialr/sns_backend/pkg/postgresql"
	"github.com/spf13/viper"
	gLog "gorm.io/gorm/logger"
)

// Http init all preparation to run the API using http server.
func Http() {
	// init viper config
	v, err := conf.InitConfigYml()
	if err != nil {
		log.Fatalln("failed to init config:", err)
	}
	// setup logger target based on the config debug state
	logFile, logGorm := setupLogger(v)
	defer logFile.Close()
	defer logGorm.Close()
	// init app logger
	appWr := logger.NewFile(logFile)
	appWr.Init()
	// enable gorm debug based on the app config
	var gormDebugLvl []gLog.LogLevel
	if !v.GetBool("server.debug") {
		// only log the error if debug is false
		gormDebugLvl = append(gormDebugLvl, gLog.Error)
	}
	// init gorm logger along with the log level defined above
	gormLog := gormLogger.New(logGorm, gormDebugLvl...)
	// init gorm using postgresql as the DB
	db, err := postgresql.NewGorm(v, gormLog)
	if err != nil {
		appWr.Err("failed to init gorm with postgresql as the DB:", err)
		return
	}
	// get the sql db
	sqlDB, err := db.DB()
	if err != nil {
		appWr.Err("failed to get the DB instance from gorm:", err)
		return
	}
	defer sqlDB.Close()
	// init fiber
	fiberApp := fiber.New(fiber.Config{
		IdleTimeout:           5 * time.Second,
		BodyLimit:             50 * 1024 * 1024, // 50MB
		RequestMethods:        []string{fiber.MethodHead, fiber.MethodGet, fiber.MethodPost},
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
		DisableStartupMessage: !v.GetBool("server.debug"),
	})
	// ini fiber log middleware with custom config
	fiberLog := fLog.New(fLog.Config{
		Output: logFile,
	})
	// add middlewares
	fiberApp.Use(
		fiberLog,
		recover.New(),
		compress.New(),
		cors.New(),
		helmet.New(),
	)
	// init http handlers
	h := app.HttpHandlers{
		R:      fiberApp.Group("/api"), // add prefix /api to route stack
		DB:     db,
		Config: v,
		Log:    appWr,
	}
	h.SetupRouter()
	// listen from a different goroutine
	go func() {
		if err := fiberApp.Listen(v.GetString("server.host") + ":" + v.GetString("server.port")); err != nil {
			log.Panicf("failed listen into port %v", err)
		}
	}()
	// create channel to signify a signal being sent
	c := make(chan os.Signal, 1)
	// when an interrupt or termination signal is sent, notify the channel
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// blocks main thread until an interrupt is received
	<-c
	log.Println("gracefully shutting down...")
	fiberApp.Shutdown()
	fmt.Println("Running cleanup tasks...")
	sqlDB.Close()
	fmt.Println("services was successful shutdown.")
}

func setupLogger(v *viper.Viper) (*os.File, *os.File) {
	logFile, logGorm, err := os.Stdout, os.Stdout, errors.New("")
	if !v.GetBool("server.debug") {
		logDir := strings.TrimSuffix(v.GetString("log.dir"), "/")
		// log file for fiber app
		targetAppLog := logDir + "/app-log"
		logFile, err = os.Open(targetAppLog)
		if err != nil {
			log.Fatalln("failed to open log file for app:", err)
		}
		// log file for gorm
		targetGormLog := logDir + "/gorm-log"
		logGorm, err = os.Open(targetGormLog)
		if err != nil {
			log.Fatalln("failed to open log file for gorm:", err)
		}
	}

	return logFile, logGorm
}
