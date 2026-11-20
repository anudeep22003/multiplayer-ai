package internal

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config holds application settings loaded from environment variables.
type Config struct {
	AppName     string `env:"APP_NAME" envDefault:"conman"`
	Environment string `env:"APP_ENV" envDefault:"development"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`

	Host         string        `env:"HOST" envDefault:"0.0.0.0"`
	Port         uint16        `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`

	DatabaseURL string `env:"DATABASE_URL"`
}

// Load reads environment variables, applying defaults for unset values.
// Invalid values return an error. DATABASE_URL is optional until needed.
func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}
	return cfg, nil
}
