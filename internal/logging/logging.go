package logging

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

func New(level string) zerolog.Logger {
	l := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	}).With().Timestamp().Logger()

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		panic("error parsing log level: " + err.Error())
	}

	zerolog.SetGlobalLevel(lvl)

	return l
}
