package logger

import (
	"io"
	"log"
)

// logFile log writer that write to a file.
type logFile struct {
	file           io.Writer
	infLog, errLog *log.Logger
}

func (l *logFile) Init() {
	l.infLog = log.New(l.file, "[INF] ", log.Ldate|log.Ltime)
	l.errLog = log.New(l.file, "[ERR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *logFile) Inf(msg ...any) {
	l.infLog.Println(msg...)
}

func (l *logFile) Err(msg ...any) {
	l.errLog.Println(msg...)
}
