package main

import (
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/config"
	"github.com/syncrepair/backend/internal/handler"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/database/mongo"
)

const configFilePath = "config.yml"

func main() {
	cfg := config.Load(configFilePath)

	mongoClient := mongo.NewClient(cfg.Mongo.URI)
	mongoDB := mongoClient.Database(cfg.Mongo.Name)

	userRepository := repository.NewUserRepository(mongoDB)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userHandler := handler.NewUserHandler(userUsecase)

	e := echo.New()

	apiGroup := e.Group("/api")
	{
		userHandler.Routes(apiGroup)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
