package main

import (
	"context"
	"github.com/syncrepair/backend/internal/config"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/server"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/database/mongo"
	"github.com/syncrepair/backend/pkg/logging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const configFilePath = "config.yml"

func main() {
	// Configuration
	cfg := config.Load(configFilePath)

	// Logging
	log := logging.New(cfg.LogLevel)

	log.Info().
		Str("version", cfg.App.Version).
		Msgf("🚀 Starting %s", cfg.App.Name)

	// Databases
	log.Info().
		Msg("Connecting to mongo database")

	mongoClient := mongo.NewClient(cfg.Mongo.URI)
	mongoDB := mongoClient.Database(cfg.Mongo.Name)

	// Repositories
	log.Info().
		Msg("Instantiating repositories")

	userRepository := repository.NewUserRepository(mongoDB)

	// Services
	log.Info().
		Msg("Instantiating usecases")

	usecase.NewUserUsecase(userRepository)

	// Server
	log.Info().
		Msg("Instantiating a server")

	srv := server.New(&server.Config{
		Port:         cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}, nil)

	log.Info().
		Int("port", cfg.Server.Port).
		Dur("read_timeout", cfg.Server.ReadTimeout).
		Dur("write_timeout", cfg.Server.WriteTimeout).
		Dur("idle_timeout", cfg.Server.IdleTimeout).
		Msg("Running the server")

	go srv.MustRun()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	log.Info().
		Msg("Stopping the server")

	srv.MustStop(ctx)

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		log.Error().
			Err(err).
			Msg("error disconnecting from mongo database")
	}
}
