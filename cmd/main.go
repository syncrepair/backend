package main

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/bootstrap/config"
	"github.com/syncrepair/backend/internal/bootstrap/postgres"
	"github.com/syncrepair/backend/internal/handler"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/usecase"
)

func main() {
	cfg := config.Init()

	postgresDB := postgres.Init(postgres.Config{
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Database: cfg.Postgres.Database,
	})
	defer postgresDB.Close()

	postgresSB := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	userRepository := repository.NewUserRepository(postgresDB, postgresSB, "users")
	userUsecase := usecase.NewUserUsecase(userRepository)
	userHandler := handler.NewUserHandler(userUsecase)

	e := echo.New()

	apiGroup := e.Group("/api")
	{
		userHandler.Routes(apiGroup)
	}

	if err := e.Start(cfg.HTTP.Addr); err != nil {
		panic("error starting server: " + err.Error())
	}
}
