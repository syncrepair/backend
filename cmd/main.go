package main

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	_ "github.com/joho/godotenv/autoload"
	"github.com/syncrepair/backend/internal/config"
	"github.com/syncrepair/backend/internal/delivery/http/handler"
	"github.com/syncrepair/backend/internal/delivery/http/server"
	"github.com/syncrepair/backend/internal/logging"
	postgresRepository "github.com/syncrepair/backend/internal/repository/postgres"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/database/postgres"
	"github.com/syncrepair/backend/pkg/password"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Configuration
	cfg := config.Load()

	// Logging
	log := logging.New(cfg.LogLevel)
	log.Info().
		Msgf("starting %s", cfg.AppName)

	// Postgres
	postgresDB := postgres.New(cfg.Postgres.URL)
	defer postgresDB.Close()

	postgresSB := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Packages
	passwordHasher := password.NewHasher(cfg.Password.Salt, cfg.Password.HashingCost)

	// Repositories
	log.Info().
		Msg("initializing postgres repositories")

	companyRepository := postgresRepository.NewCompanyRepository(postgresDB, postgresSB)
	userRepository := postgresRepository.NewUserRepository(postgresDB, postgresSB)

	// Usecases
	log.Info().
		Msg("initializing usecases")

	companyUsecase := usecase.NewCompanyUsecase(companyRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, companyRepository, passwordHasher)

	// HTTP
	log.Info().
		Msg("initializing http server")

	httpServer := server.New(server.Config{
		AppName:      cfg.AppName,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		Logger:       log,
	})

	companyHandler := handler.NewCompanyHandler(companyUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	api := httpServer.Group("/api")
	{
		companies := api.Group("/companies")
		{
			companies.Post("/", companyHandler.Create)
		}

		users := api.Group("/users")
		{
			users.Post("/signup", userHandler.SignUp)
		}
	}

	log.Info().
		Int("port", cfg.HTTPServer.Port).
		Msg("starting http server")

	go func() {
		err := httpServer.Listen(fmt.Sprintf(":%d", cfg.HTTPServer.Port))
		if err != nil {
			panic("error while starting http server: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Info().
		Msg("stopping http server")

	if err := httpServer.Shutdown(); err != nil {
		panic("error while stopping http server: " + err.Error())
	}
}
