package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/bootstrap/config"
	"github.com/syncrepair/backend/internal/handler"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/database/mongo"
)

func main() {
	cfg := config.Init()

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

	if err := e.Start(fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)); err != nil {
		panic("error starting server: " + err.Error())
	}
}
