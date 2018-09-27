package logger

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z"
	zerolog.CallerSkipFrameCount++

	log.Logger = log.With().Caller().Logger()
}

// SetLevel set log level
func SetLevel(level string) error {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}

	Infof("set log level [%s]", lvl.String())
	zerolog.SetGlobalLevel(lvl)
	return nil
}

// Debug debug log
func Debug(args ...interface{}) {
	log.Debug().Msg(fmt.Sprint(args...))
}

// Debugf debug log with format
func Debugf(format string, args ...interface{}) {
	log.Debug().Msgf(format, args...)
}

// Info info log
func Info(args ...interface{}) {
	log.Info().Msg(fmt.Sprint(args...))
}

// Infof info log with format
func Infof(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

// Warn warning log
func Warn(args ...interface{}) {
	log.Warn().Msg(fmt.Sprint(args...))
}

// Warnf warning log with format
func Warnf(format string, args ...interface{}) {
	log.Warn().Msgf(format, args...)
}

// Error error log
func Error(args ...interface{}) {
	log.Error().Msg(fmt.Sprint(args...))
}

// Errorf error log with format
func Errorf(format string, args ...interface{}) {
	log.Error().Msgf(format, args...)
}

// Panic panic log
func Panic(args ...interface{}) {
	log.Panic().Msg(fmt.Sprint(args...))
}

// Panicf panic log with format
func Panicf(format string, args ...interface{}) {
	log.Panic().Msgf(format, args...)
}
