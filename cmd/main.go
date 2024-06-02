package main

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/bootstrap/config"
	"github.com/syncrepair/backend/internal/bootstrap/logger"
	"github.com/syncrepair/backend/internal/bootstrap/postgres"
	"github.com/syncrepair/backend/internal/bootstrap/server"
	"github.com/syncrepair/backend/internal/handler"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/hasher"
	"github.com/ziflex/lecho/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	cfg := config.Init()

	log := logger.Init()
	log.Info().
		Msgf("🚀 Starting %s", cfg.App.Name)

	log.Info().
		Msg("Connecting to postgres database")

	postgresDB := postgres.Init(postgres.Config{
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Database: cfg.Postgres.Database,
	})

	if err := postgresDB.Ping(ctx); err != nil {
		log.Fatal().
			Err(err).
			Msg("error pinging postgres database")
	}
	defer postgresDB.Close()

	postgresSB := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	passwordHasher := hasher.NewPasswordHasher(cfg.Auth.PasswordSalt)

	userRepository := repository.NewUserRepository(postgresDB, postgresSB, "users")
	userUsecase := usecase.NewUserUsecase(userRepository, passwordHasher)
	userHandler := handler.NewUserHandler(userUsecase)

	h := echo.New()

	h.Logger = lecho.From(log)

	apiGroup := h.Group("/api")
	{
		userHandler.Routes(apiGroup)
	}

	log.Info().
		Str("address", cfg.HTTP.Address).
		Msg("Starting HTTP server")

	srv := server.Init(server.Config{
		Handler:      h,
		Addr:         cfg.HTTP.Address,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().
				Err(err).
				Msg("error starting http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctxTimeout, shutdown := context.WithTimeout(ctx, 5*time.Second)
	defer shutdown()

	defer func() {
		if err := srv.Shutdown(ctxTimeout); err != nil {
			log.Fatal().
				Err(err).
				Msg("error shutting down http server")
		}
	}()
}
