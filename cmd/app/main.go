package main

import (
	"context"
	"github.com/syncrepair/backend/internal/config"
	"github.com/syncrepair/backend/internal/delivery/http"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/auth"
	"github.com/syncrepair/backend/pkg/database/mongodb"
	"github.com/syncrepair/backend/pkg/http/server"
	"github.com/syncrepair/backend/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

// @title Syncrepair API
// @version 1.0
// @host localhost:80
// @BasePath /api
//
// @securityDefinitions.apikey UserAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	log := logger.New(cfg.App.Env, config.DevEnv, config.ProdEnv)

	log.Info().
		Msgf("🚀 Starting %s", cfg.App.Name)

	log.Info().
		Msg("Connecting to mongo database")

	mongoClient := mongodb.MustInit(ctx, cfg.Mongo.URI)
	mongoDB := mongoClient.Database(cfg.Mongo.Name)

	passwordHasher := auth.NewPasswordHasher(cfg.Auth.Password.Salt)
	tokensManager := auth.NewTokensManager(cfg.Auth.Tokens.AccessTokenKey, cfg.Auth.Tokens.AccessTokenTTL)

	userRepository := repository.NewUserRepository(mongoDB)
	companyRepository := repository.NewCompanyRepository(mongoDB)
	clientRepository := repository.NewClientRepository(mongoDB)
	serviceRepository := repository.NewServiceRepository(mongoDB)

	userUsecase := usecase.NewUserUsecase(userRepository, passwordHasher, tokensManager, cfg.Auth.Tokens.RefreshTokenTTL)
	companyUsecase := usecase.NewCompanyUsecase(companyRepository)
	clientUsecase := usecase.NewClientUsecase(clientRepository)
	vehicleUsecase := usecase.NewVehicleUsecase(clientRepository)
	serviceUsecase := usecase.NewServiceUsecase(serviceRepository)

	handler := http.NewHandler(log, tokensManager, http.Usecases{
		User:    userUsecase,
		Company: companyUsecase,
		Service: serviceUsecase,
		Client:  clientUsecase,
		Vehicle: vehicleUsecase,
	})

	srv := server.New(handler.Init(), server.Config{
		Addr:         cfg.HTTP.Address,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	})

	log.Info().
		Msg("Starting HTTP server")

	go func() {
		if err := srv.Run(); err != nil {
			log.Error().Msgf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctxTimeout, shutdown := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer shutdown()

	if err := srv.Stop(ctxTimeout); err != nil {
		log.Fatal().
			Err(err).
			Msg("error occurred while shutting down http server")
	}

	if err := mongoClient.Disconnect(ctxTimeout); err != nil {
		log.Fatal().
			Err(err).
			Msg("error disconnecting from mongo database")
	}
}
