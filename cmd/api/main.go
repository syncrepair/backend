package main

import (
	"github.com/syncrepair/backend/internal/config"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/database/mongo"
	"github.com/syncrepair/backend/pkg/logging"
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

	mongoDB := mongo.NewClient(cfg.Mongo.URI).Database(cfg.Mongo.Name)

	// Repositories
	log.Info().
		Msg("Instantiating repositories")

	userRepository := repository.NewUserRepository(mongoDB)

	// Services
	log.Info().
		Msg("Instantiating usecases")

	usecase.NewUserUsecase(userRepository)
}
