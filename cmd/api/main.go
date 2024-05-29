package main

import (
	"github.com/syncrepair/backend/config"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/service"
	"github.com/syncrepair/backend/pkg/database/mongo"
	"github.com/syncrepair/backend/pkg/logging"
	"os"
)

func main() {
	// Configuration
	cfg := config.Init()

	// Logger
	logger := logging.New(os.Stderr, cfg.LogLevel)

	logger.Info().
		Msgf("🚀 Starting %s", cfg.AppName)

	// Database
	logger.Info().
		Msg("Connecting to mongo database")

	mongoDB := mongo.NewClient(cfg.MongoURI).Database(cfg.MongoName)

	// Repositories
	userRepository := repository.NewUserRepository(mongoDB.Collection("users"))

	// Services
	service.NewUserService(userRepository)
}
