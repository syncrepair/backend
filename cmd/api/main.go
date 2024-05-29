package main

import (
	"github.com/syncrepair/backend/config"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/service"
	"github.com/syncrepair/backend/pkg/database/mongo"
)

const configFilePath = "config.yml"

func main() {
	// Configuration
	cfg := config.Load(configFilePath)

	// Database
	mongoDB := mongo.NewClient(cfg.Mongo.URI).Database(cfg.Mongo.Name)

	// Repositories
	userRepository := repository.NewUserRepository(mongoDB.Collection("users"))

	// Services
	service.NewUserService(userRepository)
}
