package logging

import (
	"io"

	"github.com/sirupsen/logrus"
)

type LogrusAdpter struct {
	logger *logrus.Logger
}

func NewLogrusAdpter(logger *logrus.Logger) *LogrusAdpter {
	return &LogrusAdpter{logger}
}

func (l *LogrusAdpter) Info(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *LogrusAdpter) Debug(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *LogrusAdpter) Error(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *LogrusAdpter) Warn(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *LogrusAdpter) Fatal(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *LogrusAdpter) SetOutput(w io.Writer) {
	l.logger.Out = w
}

func (l *LogrusAdpter) Println(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *LogrusAdpter) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}
