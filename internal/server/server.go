package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	srv *http.Server
}

type Config struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func New(cfg *Config, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) MustStop(ctx context.Context) {
	if err := s.Stop(ctx); err != nil {
		log.Fatalf("error stopping server: %v", err)
	}
}
