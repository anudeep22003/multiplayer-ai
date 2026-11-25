package internal

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config holds application settings loaded from environment variables.
type Config struct {
	AppName     string `env:"APP_NAME" envDefault:"multiplayer-ai"`
	Environment string `env:"APP_ENV" envDefault:"development"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`

	Host         string        `env:"HOST" envDefault:"0.0.0.0"`
	Port         uint16        `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`

	DatabaseURL       string        `env:"DATABASE_URL,notEmpty"`
	DBMaxConns        int32         `env:"DB_MAX_CONNS" envDefault:"25"`
	DBMinConns        int32         `env:"DB_MIN_CONNS" envDefault:"2"`
	DBMaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME" envDefault:"1h"`
	DBMaxConnIdleTime time.Duration `env:"DB_MAX_CONN_IDLE_TIME" envDefault:"30m"`
	DBConnectTimeout  time.Duration `env:"DB_CONNECT_TIMEOUT" envDefault:"5s"`
}

// Load reads environment variables, applying defaults for unset values.
// Invalid or missing required values return an error.
func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}
	return cfg, nil
}
