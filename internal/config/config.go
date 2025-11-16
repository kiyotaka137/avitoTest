package config

import (
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	Http     HttpConfig
	Database DatabaseConfig
	Log      LogConfig
	JWT      JWTConfig
}

type HttpConfig struct {
	Address     string        `env:"HTTP_ADDRESS" envDefault:"localhost:8080"`
	Timeout     time.Duration `env:"HTTP_TIMEOUT" envDefault:"10s"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
}

type JWTConfig struct {
	Secret          string        `env:"SECRET" envDefault:"secret"`
	AccessTokenTTL  time.Duration `env:"ACCESS_TTL" envDefault:"15m"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TTL" envDefault:"24h"`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	DBName   string `env:"DB_NAME" envDefault:"postgres"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
	Password string `env:"DB_PASSWORD" envDefault:""`
	User     string `env:"DB_USER" envDefault:"postgres"`
}

type LogConfig struct {
	Level string `env:"LEVEL" envDefault:"info"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	opts := env.Options{
		UseFieldNameByDefault: false,
		Environment:           nil,
		TagName:               "env",
	}

	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, err
	}
	return cfg, nil
}
