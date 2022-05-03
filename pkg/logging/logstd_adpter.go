package logging

import (
	"io"
	"log"
)

type LogStdAdapter struct {
	logger *log.Logger
}

func NewLogStdAdapter() *LogStdAdapter {
	return &LogStdAdapter{logger: log.New(log.Writer(), "", log.LstdFlags)}
}

func (l *LogStdAdapter) Info(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *LogStdAdapter) Debug(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *LogStdAdapter) Error(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *LogStdAdapter) Warn(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *LogStdAdapter) Fatal(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *LogStdAdapter) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *LogStdAdapter) Println(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *LogStdAdapter) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}
