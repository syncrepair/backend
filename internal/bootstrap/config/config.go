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

	cfg.HTTP.Addr = os.Getenv("HTTP_ADDR")

	cfg.Postgres.Username = os.Getenv("POSTGRES_USERNAME")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.Database = os.Getenv("POSTGRES_DATABASE")
}
