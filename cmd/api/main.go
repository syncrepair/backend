package main

import (
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/config"
	"github.com/syncrepair/backend/internal/handler"
	"github.com/syncrepair/backend/internal/logger"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/database/mongo"
)

const configFilePath = "config.yml"

func main() {
	// Configuration
	cfg := config.Load(configFilePath)

	// Logging
	l := logger.New(cfg.LogLevel)

	l.Info().
		Str("version", cfg.App.Version).
		Msgf("🚀 Starting %s", cfg.App.Name)

	// Databases
	l.Info().
		Msg("Connecting to mongo database")

	mongoClient := mongo.NewClient(cfg.Mongo.URI)
	mongoDB := mongoClient.Database(cfg.Mongo.Name)

	// Repositories
	l.Info().
		Msg("Instantiating repositories")

	userRepository := repository.NewUserRepository(mongoDB)

	// Services
	l.Info().
		Msg("Instantiating usecases")

	userUsecase := usecase.NewUserUsecase(userRepository)

	// Handlers
	l.Info().
		Msg("Instantiating handlers")

	userHandler := handler.NewUserHandler(userUsecase)

	// Router
	l.Info().
		Msg("Instantiating router")

	e := echo.New()

	apiGroup := e.Group("/api")
	{
		userHandler.Routes(apiGroup)
	}

	l.Info().
		Msg("Starting server")

	e.Logger.Fatal(e.Start(":8080"))
}
