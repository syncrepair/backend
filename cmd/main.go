package main

import (
	"context"
	"github.com/syncrepair/backend/config"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/pkg/db/mongo"
)

func main() {
	ctx := context.Background()

	cfg := config.Init()

	mongoClient := mongo.NewClient(ctx, cfg.MongoURI)
	defer mongoClient.Disconnect(ctx)

	mongoDB := mongoClient.Database(cfg.MongoName)

	repository.NewUser(mongoDB.Collection("users"))
}
