package logger

import (
	"github.com/rs/zerolog"
	"log"
	"os"
	"time"
)

type Logger struct {
	zerolog.Logger
}

func New(level string) *Logger {
	l := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
	}).With().Timestamp().Logger()

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Fatalf("error parsing log level: %v", err)
	}
	l = l.Level(lvl)

	return &Logger{l}
}
