package logger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

func Init() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
	}).With().Timestamp().Logger()
}
