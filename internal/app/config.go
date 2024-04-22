package app

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	Database struct {
		Host     string `env:"DATABASE_HOST" envDefault:"127.0.0.1"`
		Port     string `env:"DATABASE_PORT" envDefault:"3307"`
		User     string `env:"DATABASE_USER" envDefault:"user"`
		Password string `env:"DATABASE_PASSWORD" envDefault:"password"`
		Name     string `env:"DATABASE_NAME" envDefault:"norddb"`
	}
	App struct {
		Port       string `env:"APP_PORT" envDefault:":8080"`
		TimeFormat string `env:"APP_TIME_FORMAT" envDefault:"2006-01-02T15:04:05"`
	}
}

func LoadConfig() (*Config, error) {
	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
