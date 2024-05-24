package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	AppName string `env:"APP_NAME" env-default:"backend"`
}

var (
	cfg  *Config
	once sync.Once
)

func Init() *Config {
	once.Do(func() {
		cfg = &Config{}

		if err := cleanenv.ReadConfig(".env", cfg); err != nil {
			log.Fatalf("error reading .env file: %v", err)
		}
	})

	return cfg
}
