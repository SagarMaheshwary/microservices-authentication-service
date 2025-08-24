package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "02/01/2006, 3:04:05 PM",
	})
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

func Panic(format string, v ...interface{}) {
	log.Panic().Msgf(format, v...)
}
