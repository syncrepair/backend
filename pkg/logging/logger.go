package logging

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

type Logger struct {
	zerolog.Logger
}

func New() *Logger {
	l := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
	}).With().Timestamp().Logger()

	return &Logger{l}
}
