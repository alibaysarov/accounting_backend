package settings

import (
	"errors"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	DbUrl   string `envconfig:"DB_URL"`
	Port    string `envconfig:"APP_PORT"`
	GinMode string `envconfig:"GIN_MODE" default:"release"`

	JwtKey string `envconfig:"JWT_SECRET_KEY"`
}

func (cfg *AppConfig) Init() error {
	if err := godotenv.Overload(); err != nil {
		return errors.New("no .env file found, using OS env vars")
	}

	return envconfig.Process("", cfg)
}
