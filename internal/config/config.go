package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"sync"
)

type (
	Config struct {
		Env      string
		LogLevel string `yaml:"log_level" env-required:"true"`
	}
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

var config Config
var once sync.Once

func Load() *Config {
	once.Do(func() {
		if os.Getenv("ENV") == EnvDev || os.Getenv("ENV") == EnvProd {
			config.Env = os.Getenv("ENV")
		} else {
			panic("ENV is invalid")
		}

		switch config.Env {
		case EnvDev:
			if err := cleanenv.ReadConfig("configs/dev.yml", &config); err != nil {
				panic("error reading config: " + err.Error())
			}
		case EnvProd:
			if err := cleanenv.ReadConfig("configs/prod.yml", &config); err != nil {
				panic("error reading config: " + err.Error())
			}
		}
	})

	return &config
}
