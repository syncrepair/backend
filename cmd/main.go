package main

import (
	"context"
	"github.com/syncrepair/backend/config"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/pkg/db/mongo"
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
	mongoClient := mongo.NewClient(ctx, cfg.MongoURI)
	defer mongoClient.Disconnect(ctx)

	mongoDB := mongoClient.Database(cfg.MongoName)

	// Repositories
	repository.NewUser(mongoDB.Collection("users"))
}
