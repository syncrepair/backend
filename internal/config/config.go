package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type (
	Config struct {
		LogLevel string      `yaml:"log_level" env-default:"error"`
		App      AppConfig   `yaml:"app" env-required:"true"`
		Mongo    MongoConfig `yaml:"mongo" env-required:"true"`
	}

	AppConfig struct {
		Name    string `yaml:"name" env-required:"true"`
		Version string `yaml:"version" env-required:"true"`
	}

	MongoConfig struct {
		URI  string `yaml:"uri" env-required:"true"`
		Name string `yaml:"name" env-required:"true"`
	}
)

var (
	cfg  *Config
	once sync.Once
)

func Load(filePath string) *Config {
	once.Do(func() {
		cfg = &Config{}

		if err := cleanenv.ReadConfig(filePath, cfg); err != nil {
			log.Fatalf("error reading %s file: %v", filePath, err)
		}
	})

	return cfg
}
