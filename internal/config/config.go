package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"sync"
	"time"
)

type (
	Config struct {
		Env        string
		AppName    string     `yaml:"app_name" env-required:"true"`
		LogLevel   string     `yaml:"log_level" env-required:"true"`
		HTTPServer HTTPServer `yaml:"http_server" env-required:"true"`
		Postgres   Postgres   `yaml:"postgres" env-required:"true"`
		Password   Password   `yaml:"password" env-required:"true"`
	}

	HTTPServer struct {
		Port         int           `yaml:"port" env-required:"true"`
		ReadTimeout  time.Duration `yaml:"read_timeout" env-required:"true"`
		WriteTimeout time.Duration `yaml:"write_timeout" env-required:"true"`
		IdleTimeout  time.Duration `yaml:"idle_timeout" env-required:"true"`
	}

	Postgres struct {
		URL string `yaml:"url" env:"POSTGRES_URL" env-required:"true"`
	}

	Password struct {
		HashingCost int    `yaml:"hashing_cost" env-required:"true"`
		Salt        string `yaml:"salt" env:"PASSWORD_SALT" env-required:"true"`
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
