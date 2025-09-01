package config

import (
	"errors"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
)

type TgConfig struct {
	Bot BotConfig
	DB  DatabaseConfig
}

type BotConfig struct {
	Token      string `env:"TELEGRAM_TOKEN"`
	BotTimeout int    `env:"TELEGRAM_TIMEOUT"     envDefault:"10"`
	BotOffset  int    `env:"TELEGRAM_OFFSET"     envDefault:"0"`
}

type DatabaseConfig struct {
	Host              string        `env:"DATABASE_HOST"`
	Port              string        `env:"DATABASE_PORT"`
	DBName            string        `env:"DATABASE_NAME"`
	DBUser            string        `env:"DATABASE_USER"`
	DBPassword        string        `env:"DATABASE_PASSWORD"`
	DBSSLMode         string        `env:"DATABASE_SSL_MODE"          envDefault:"disable"`
	DBMaxOpenConn     int           `env:"DATABASE_MAX_OPEN_CONN"     envDefault:"0"`
	DBMaxIdleConn     int           `env:"DATABASE_MAX_IDLE_CONN"     envDefault:"0"`
	DBConnMaxLifetime time.Duration `env:"DATABASE_CONN_MAX_LIFETIME" envDefault:"10s"`
}

func MustNew() *TgConfig {
	cfgEnv := TgConfig{}

	err := env.Parse(&cfgEnv)
	if err != nil {
		log.Fatalf("err loading file: %s", err)
	}

	errs := cfgEnv.Validate()
	if errs != nil {
		log.Fatalf("err validating config: %s", errs.Error())
	}

	return &cfgEnv
}

func (cfg *TgConfig) Validate() (result error) {
	if cfg.Bot.Token == "" {
		result = errors.Join(result, ErrNoBotToken)
	}

	if cfg.DB.Host == "" {
		result = errors.Join(result, ErrNoDBHost)
	}

	if cfg.DB.Port == "" {
		result = errors.Join(result, ErrNoDBPort)
	}

	if cfg.DB.DBName == "" {
		result = errors.Join(result, ErrNoDBName)
	}

	if cfg.DB.DBUser == "" {
		result = errors.Join(result, ErrNoDBUser)
	}

	if cfg.DB.DBPassword == "" {
		result = errors.Join(result, ErrNoDBPassword)
	}

	return result
}
