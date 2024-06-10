package http

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerConfig struct {
	Addr         string
	Handler      http.Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Logger       zerolog.Logger
}

type Server struct {
	*http.Server
}

var log zerolog.Logger

func NewServer(cfg ServerConfig) *Server {
	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      cfg.Handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log = cfg.Logger

	return &Server{srv}
}

func (srv *Server) StartWithGracefulShutdown(ctx context.Context) {
	log.Info().
		Str("address", srv.Addr).
		Msg("Starting HTTP server")

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().
				Err(err).
				Msg("error starting http http")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctxTimeout, shutdown := context.WithTimeout(ctx, 5*time.Second)
	defer shutdown()

	log.Info().
		Msg("Gracefully shutting down HTTP server")

	defer func() {
		if err := srv.Shutdown(ctxTimeout); err != nil {
			log.Fatal().
				Err(err).
				Msg("error shutting down http server")
		}
	}()
}
