package settings

import (
	"errors"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	DbUrl string `envconfig:"DB_URL"`
	Port  string `envconfig:"APP_PORT"`
}

func (cfg *AppConfig) Init() error {
	if err := godotenv.Overload(); err != nil {
		return errors.New("no .env file found, using OS env vars")
	}

	return envconfig.Process("", cfg)
}
