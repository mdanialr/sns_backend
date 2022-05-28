package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/mdanialr/sns_backend/internal/service"
)

var (
	InfL *log.Logger
	ErrL *log.Logger
)

// InitLogger init and setup log file to write internal app log.
func InitLogger(conf *service.Config) error {
	fl, err := os.OpenFile(conf.LogDir+"sns_backend-internal-log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		return fmt.Errorf("failed to open|create log file: %v", err)
	}

	InfL = log.New(fl, "[INFO] ", log.Ldate|log.Ltime)
	ErrL = log.New(fl, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
