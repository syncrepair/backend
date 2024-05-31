package config

import "time"

type Config struct {
	App struct {
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
}
