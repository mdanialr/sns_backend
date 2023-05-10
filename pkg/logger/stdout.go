package logger

import (
	"log"
	"os"
)

// logStdOut log writer that write to a file.
type logStdOut struct {
	infLog, errLog *log.Logger
}

func (l *logStdOut) Init() {
	l.infLog = log.New(os.Stdout, "[INF] ", log.Ldate|log.Ltime)
	l.errLog = log.New(os.Stdout, "[ERR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *logStdOut) Inf(msg ...any) {
	l.infLog.Println(msg...)
}

func (l *logStdOut) Err(msg ...any) {
	l.errLog.Println(msg...)
}
