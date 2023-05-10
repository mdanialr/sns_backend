package logger

import "io"

// Writer interface to write log message.
type Writer interface {
	// Init do some necessary setup before calling WriteInf & WriteErr, Ideally
	// this only called once.
	Init()
	// Inf write message for info-level log
	Inf(...any)
	// Err write message for error-level log
	Err(...any)
}

// NewFile return logger that write message to given io.Writer.
func NewFile(wr io.Writer) Writer {
	return &logFile{file: wr}
}

// NewStdOut return logger that write message to os.Stdout.
func NewStdOut() Writer {
	return &logStdOut{}
}
