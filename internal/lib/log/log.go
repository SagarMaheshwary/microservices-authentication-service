package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func Info(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func Warn(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}

func Debug(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

func Error(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}
