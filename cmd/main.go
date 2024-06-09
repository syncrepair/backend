package main

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/syncrepair/backend/internal/config"
	"github.com/syncrepair/backend/internal/controller"
	"github.com/syncrepair/backend/internal/logger"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/auth"
	"github.com/syncrepair/backend/pkg/postgres"
	"github.com/syncrepair/backend/pkg/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	cfg := config.Init()

	log := logger.Init()
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

	log.Info().
		Str("address", cfg.HTTP.Address).
		Msg("Starting HTTP server")

	srv := server.Init(server.Config{
		Handler:      r,
		Addr:         cfg.HTTP.Address,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().
				Err(err).
				Msg("error starting http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctxTimeout, shutdown := context.WithTimeout(ctx, 5*time.Second)
	defer shutdown()

	defer func() {
		if err := srv.Shutdown(ctxTimeout); err != nil {
			log.Fatal().
				Err(err).
				Msg("error shutting down http server")
		}
	}()
}
