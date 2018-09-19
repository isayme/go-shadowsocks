package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var log = logrus.New()

func init() {
	log.Formatter = &prefixed.TextFormatter{
		FullTimestamp: true,
	}
	log.Level = logrus.DebugLevel
}

// Debug debug log
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Debugf debug log with format
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Info info log
func Info(args ...interface{}) {
	log.Info(args...)
}

// Infof info log with format
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warn warning log
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Warnf warning log with format
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Error error log
func Error(args ...interface{}) {
	log.Error(args...)
}

// Errorf error log with format
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Panic panic log
func Panic(args ...interface{}) {
	log.Panic(args...)
}

// Panicf panic log with format
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// Println print log
func Println(args ...interface{}) {
	log.Println(args...)
}

// Printlnf print log with format
func Printlnf(format string, args ...interface{}) {
	log.Println(fmt.Sprintf(format, args...))
}
