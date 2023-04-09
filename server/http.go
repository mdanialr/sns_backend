package server

import (
	"errors"
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
	"github.com/mdanialr/sns_backend/pkg/helper"
	"github.com/mdanialr/sns_backend/pkg/logger"
	"github.com/mdanialr/sns_backend/pkg/postgresql"
	"github.com/mdanialr/sns_backend/pkg/storage"
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
	logFile, logGorm, logFiber := setupLogger(v)
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
		os.Exit(1)
		return
	}
	// get the sql db
	sqlDB, err := db.DB()
	if err != nil {
		appWr.Err("failed to get the DB instance from gorm:", err)
		os.Exit(1)
		return
	}
	defer sqlDB.Close()
	// init storage
	var st storage.IStorage
	switch v.GetString("storage.driver") {
	case "file":
		st = storage.NewFile(appWr)
	default:
		appWr.Err("unsupported storage driver. currently support [file]")
		os.Exit(1)
		return
	}
	// init fiber
	fiberApp := fiber.New(fiber.Config{
		IdleTimeout:           5 * time.Second,
		BodyLimit:             v.GetInt("server.limit") * 1024 * 1024,
		RequestMethods:        []string{fiber.MethodHead, fiber.MethodGet, fiber.MethodPost},
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
		DisableStartupMessage: !v.GetBool("server.debug"),
		ErrorHandler:          helper.DefaultHTTPErrorHandler,
	})
	// ini fiber log middleware with custom config
	fiberLog := fLog.New(fLog.Config{
		Output:     logFiber,
		TimeFormat: "02-Jan-06 15:04:05",
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
		R:       fiberApp.Group("/api"), // add prefix /api to route stack
		DB:      db,
		Config:  v,
		Log:     appWr,
		Storage: st,
	}
	h.SetupRouter()
	// log the app host and port
	host := v.GetString("server.host") + ":" + v.GetString("server.port")
	appWr.Inf("Run app in", host)
	// listen from a different goroutine
	go func() {
		if err := fiberApp.Listen(host); err != nil {
			log.Panicf("failed listen into port %v", err)
		}
	}()
	// create channel to signify a signal being sent
	c := make(chan os.Signal, 1)
	// when an interrupt or termination signal is sent, notify the channel
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// blocks main thread until an interrupt is received
	<-c
	appWr.Inf("gracefully shutting down...")
	fiberApp.Shutdown()
	appWr.Inf("running cleanup tasks...")
	sqlDB.Close()
	appWr.Inf("services was successful shutdown.")
}

// setupLogger return logger target file for internal log, gorm log, and fiber
// app log
func setupLogger(v *viper.Viper) (*os.File, *os.File, *os.File) {
	logFile, logGorm, logFiber := os.Stdout, os.Stdout, os.Stdout
	if !v.GetBool("server.debug") {
		logDir := strings.TrimSuffix(v.GetString("log.dir"), "/")
		// log file for fiber app
		logFiber = setupLogFiber(logDir)
		// log file for internal app
		logFile = setupLog(logDir)
		// log file for gorm
		logGorm = setupLogGorm(logDir)
	}

	return logFile, logGorm, logFiber
}

func setupLogFiber(logDir string) *os.File {
	logFiber, err := os.Stdout, errors.New("")
	// log file for internal app
	targetLog := logDir + "/app-log"
	logFiber, err = os.OpenFile(targetLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln("failed to open log file for fiber app:", err)
	}
	return logFiber
}

func setupLog(logDir string) *os.File {
	logFile, err := os.Stdout, errors.New("")
	// log file for internal app
	targetLog := logDir + "/log"
	logFile, err = os.OpenFile(targetLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln("failed to open log file for internal:", err)
	}
	return logFile
}

func setupLogGorm(logDir string) *os.File {
	logGorm, err := os.Stdout, errors.New("")
	// log file for internal app
	targetLog := logDir + "/gorm-log"
	logGorm, err = os.OpenFile(targetLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln("failed to open log file for gorm:", err)
	}
	return logGorm
}
