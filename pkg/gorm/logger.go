package gorm_logger

import (
	"io"
	"log"
	"time"

	"gorm.io/gorm/logger"
)

// New use given writer as the target where the log for GORM is written to.
// Also set the log level to Info/Debug by default if logLevels is not
// provided.
func New(wr io.Writer, logLevels ...logger.LogLevel) logger.Interface {
	logLvl := logger.Info
	if len(logLevels) > 0 {
		logLvl = logLevels[0]
	}
	newGormLogger := logger.New(
		log.New(wr, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLvl,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	return newGormLogger
}
