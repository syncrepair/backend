package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"time"
)

const (
	DevEnv  = "dev"
	ProdEnv = "prod"
)

type Config struct {
	App struct {
		Env  string `yaml:"env" env:"APP_ENV" env-required:"true"`
		Name string `yaml:"name" env:"APP_NAME" env-required:"true"`
	} `yaml:"app" env-required:"true"`

	HTTP struct {
		Address      string        `yaml:"address" env:"HTTP_ADDRESS" env-required:"true"`
		ReadTimeout  time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-required:"true"`
		WriteTimeout time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-required:"true"`
		IdleTimeout  time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-required:"true"`
	} `yaml:"http" env-required:"true"`

	Postgres struct {
		Username string `yaml:"username" env:"POSTGRES_USERNAME" env-required:"true"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
		Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
		Port     int    `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
		Database string `yaml:"database" env:"POSTGRES_DATABASE" env-required:"true"`
	} `yaml:"postgres" env-required:"true"`

	Auth struct {
		JWT struct {
			Key string        `yaml:"key" env:"AUTH_JWT_KEY" env-required:"true"`
			TTL time.Duration `yaml:"ttl" env:"AUTH_JWT_TTL" env-required:"true"`
		} `yaml:"jwt" env-required:"true"`

		PasswordSalt string `yaml:"password_salt" env:"AUTH_PASSWORD_SALT" env-required:"true"`
	} `yaml:"auth" env-required:"true"`
}

func Load() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		panic("error reading config file: " + err.Error())
	}

	return &cfg
}
