package logging

import (
	"github.com/rs/zerolog"
	"io"
	"time"
)

func New(out io.Writer, level string) zerolog.Logger {
	l := zerolog.New(zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: time.DateTime,
	}).With().Timestamp().Logger()

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		panic("error parsing log level: " + err.Error())
	}

	zerolog.SetGlobalLevel(lvl)

	return l
}
