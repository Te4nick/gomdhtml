package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func Setup(debug bool) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	return Logger
}

func Get() zerolog.Logger {
	return Logger
}

func Debug(msg string) {
	Logger.Debug().Msg(msg)
}

func Debugf(format string, v ...interface{}) {
	Logger.Debug().Msgf(format, v...)
}

func Info(msg string) {
	Logger.Info().Msg(msg)
}

func Infof(format string, v ...interface{}) {
	Logger.Info().Msgf(format, v...)
}

func Warn(msg string) {
	Logger.Warn().Msg(msg)
}

func Warnf(format string, v ...interface{}) {
	Logger.Warn().Msgf(format, v...)
}

func Err(err error, msg string) {
	Logger.Err(err).Msg(msg)
}

func Errf(err error, format string, v ...interface{}) {
	Logger.Err(err).Msgf(format, v...)
}
