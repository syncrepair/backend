package logger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

func New(env, devEnv, prodEnv string) zerolog.Logger {
	var l zerolog.Logger

	switch env {
	case devEnv:
		l = devEnvLogger()
	case prodEnv:
		l = prodEnvLogger()
	default:
		l = prodEnvLogger()
	}

	return l
}

func devEnvLogger() zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	return zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
	}).With().Timestamp().Logger()
}

func prodEnvLogger() zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	return zerolog.New(os.Stderr).
		With().Timestamp().Logger()
}
