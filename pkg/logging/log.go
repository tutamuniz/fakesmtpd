package logging

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Error(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Println(args ...interface{})
	Printf(format string, args ...interface{})
	SetOutput(io.Writer)
}

const (
	// LogLevelDebug is the debug log level
	LogLevelDebug = "debug"
	// LogLevelInfo is the info log level
	LogLevelInfo = "info"
	// LogLevelWarn is the warn log level
	LogLevelWarn = "warn"
	// LogLevelError is the error log level
	LogLevelError = "error"
	// LogLevelFatal is the fatal log level
	LogLevelFatal = "fatal"
)

// LogLevel is the current log level
var LogLevel = LogLevelInfo

type Logging struct {
	Logger
}

func NewLogrusLogging() *Logging {
	logger := logrus.New()

	return &Logging{
		Logger: NewLogrusAdpter(logger),
	}
}

func NewLogStdLogging() *Logging {
	return &Logging{
		Logger: NewLogStdAdapter(),
	}
}
