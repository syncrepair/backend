package main

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/syncrepair/backend/internal/config"
	"github.com/syncrepair/backend/internal/controller"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/auth"
	"github.com/syncrepair/backend/pkg/http"
	"github.com/syncrepair/backend/pkg/logger"
	"github.com/syncrepair/backend/pkg/postgres"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	log := logger.New(cfg.App.Env, config.DevEnv, config.ProdEnv)

	log.Info().
		Msgf("🚀 Starting %s", cfg.App.Name)

	log.Info().
		Msg("Connecting to postgres database")

	postgresDB := postgres.New(postgres.Config{
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Database: cfg.Postgres.Database,
	})

	if err := postgresDB.Ping(ctx); err != nil {
		log.Fatal().
			Err(err).
			Msg("error pinging postgres database")
	}
	defer postgresDB.Close()

	postgresSB := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	passwordHasher := auth.NewPasswordHasher(cfg.Auth.PasswordSalt)
	jwtManager := auth.NewJWTManager(cfg.Auth.JWT.Key, cfg.Auth.JWT.TTL)

	userRepository := repository.NewUserRepository(nil, postgresSB, "users")
	userUsecase := usecase.NewUserUsecase(userRepository, passwordHasher, jwtManager)
	userController := controller.NewUserController(userUsecase)
	companyRepository := repository.NewCompanyRepository(nil, postgresSB, "companies")
	companyUsecase := usecase.NewCompanyUsecase(companyRepository)
	companyController := controller.NewCompanyController(companyUsecase)

	r := controller.NewRouter(log)

	apiGroup := r.Group("/api")
	{
		userController.InitRoutes(apiGroup)
		companyController.InitRoutes(apiGroup)
	}

	httpServer := http.NewServer(http.ServerConfig{
		Handler:      r,
		Addr:         cfg.HTTP.Address,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
		Logger:       log,
	})

	httpServer.StartWithGracefulShutdown(ctx)
}
