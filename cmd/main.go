package main

import (
	"context"
	"github.com/syncrepair/backend/config"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/service"
	"github.com/syncrepair/backend/pkg/database/mongo"
	"github.com/syncrepair/backend/pkg/logging"
	"os"
)

func main() {
	// Context
	ctx := context.Background()

	// Configuration
	cfg := config.Init()

	// Logger
	logger := logging.New(os.Stderr, cfg.LogLevel)

	logger.Info().
		Msgf("🚀 Starting %s", cfg.AppName)

	// Database
	logger.Info().
		Msg("Connecting to mongo database")

	mongoClient := mongo.NewClient(ctx, cfg.MongoURI)

	mongoDatabase := mongoClient.Database(cfg.MongoName)

	// Repositories
	userRepository := repository.NewUserRepository(mongoDatabase.Collection("users"))

	// Services
	service.NewUserService(userRepository)
}
