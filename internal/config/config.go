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

	Redis struct {
		URL string `yaml:"url" env:"REDIS_URL" env-required:"true"`
	} `yaml:"redis" env-required:"true"`

	Auth struct {
		Tokens struct {
			AccessTokenKey  string        `yaml:"access_token_key" env:"AUTH_TOKENS_ACCESS_TOKEN_KEY" env-required:"true"`
			AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"AUTH_TOKENS_ACCESS_TOKEN_TTL" env-required:"true"`
			RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"AUTH_TOKENS_REFRESH_TOKEN_TTL" env-required:"true"`
		} `yaml:"tokens" env-required:"true"`

		Password struct {
			Salt string `yaml:"salt" env:"AUTH_PASSWORD_SALT" env-required:"true"`
		} `yaml:"password" env-required:"true"`
	} `yaml:"auth" env-required:"true"`
}

func Load() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		panic("error reading config file: " + err.Error())
	}

	return &cfg
}
