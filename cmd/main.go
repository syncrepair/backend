package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/syncrepair/backend/internal/config"
	"github.com/syncrepair/backend/internal/delivery/http/handler"
	"github.com/syncrepair/backend/internal/delivery/http/server"
	"github.com/syncrepair/backend/internal/repository/postgres"
	"github.com/syncrepair/backend/internal/usecase"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Configuration
	cfg := config.Load()

	// Repositories
	companyRepository := postgres.NewCompanyRepository()

	// Usecases
	companyUsecase := usecase.NewCompanyUsecase(companyRepository)

	// HTTP
	httpServer := server.New(server.Config{
		AppName:      cfg.AppName,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	})

	companyHandler := handler.NewCompanyHandler(companyUsecase)

	api := httpServer.Group("/api")
	{
		company := api.Group("/companies")
		{
			company.Post("/", companyHandler.Create)
		}
	}

	go func() {
		err := httpServer.Listen(fmt.Sprintf(":%d", cfg.HTTPServer.Port))
		if err != nil {
			panic("error while starting http server: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err := httpServer.Shutdown(); err != nil {
		panic("error while stopping http server: " + err.Error())
	}
}
