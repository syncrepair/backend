package config

import (
	"github.com/joho/godotenv"
	"os"
)

func Init() *Config {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		panic("error loading .env file: " + err.Error())
	}

	setValsFromEnv(&cfg)

	return &cfg
}

func setValsFromEnv(cfg *Config) {
	cfg.App.Name = os.Getenv("APP_NAME")

	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	cfg.HTTP.Port = os.Getenv("HTTP_PORT")

	cfg.Mongo.URI = os.Getenv("MONGO_URI")
	cfg.Mongo.Name = os.Getenv("MONGO_NAME")
}
