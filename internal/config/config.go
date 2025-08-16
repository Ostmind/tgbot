package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/hashicorp/go-multierror"
)

type AppConfig struct {
	Srv  ServerConfig
	DB   DatabaseConfig
	Auth AuthConfig
	Log  LoggingConfig
}

type ServerConfig struct {
	Host                  string        `env:"SERVER_HOST"`
	Port                  int           `env:"SERVER_PORT"`
	ServerReadTimeout     time.Duration `env:"SERVER_READ_TIMEOUT"     envDefault:"5s"`
	ServerWriteTimeout    time.Duration `env:"SERVER_WRITE_TIMEOUT"    envDefault:"5s"`
	ServerIdleTimeout     time.Duration `env:"SERVER_IDLE_TIMEOUT"     envDefault:"10s"`
	ServerShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" envDefault:"10s"`
	EnvType               string        `env:"ENV_TYPE"                envDefault:"local"`
	MigrationPath         string        `env:"MIGRATION_PATH"          envDefault:"./internal/migrations"`
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

type AuthConfig struct {
	JWTSecret          string        `env:"JWT_SECRET"`
	JWTAccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN"  envDefault:"2h"`
	JWTRefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN" envDefault:"48h"`
	BcryptCost         int           `env:"BCRYPT_COST"       envDefault:"6"`
}

type LoggingConfig struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

func MustNew() *AppConfig {
	cfgEnv := AppConfig{}

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

func (cfg *AppConfig) Validate() error {
	var result *multierror.Error

	if cfg.Srv.Host == "" {
		result = multierror.Append(result, ErrNoServerHost)
	}

	if cfg.Srv.Port == 0 {
		result = multierror.Append(result, ErrNoServerPort)
	}

	if cfg.DB.Host == "" {
		result = multierror.Append(result, ErrNoDBHost)
	}

	if cfg.DB.Port == "" {
		result = multierror.Append(result, ErrNoDBPort)
	}

	if cfg.DB.DBName == "" {
		result = multierror.Append(result, ErrNoDBName)
	}

	if cfg.DB.DBUser == "" {
		result = multierror.Append(result, ErrNoDBUser)
	}

	if cfg.DB.DBPassword == "" {
		result = multierror.Append(result, ErrNoDBPassword)
	}

	if cfg.Auth.JWTSecret == "" {
		result = multierror.Append(result, ErrNoJWTSecret)
	}

	return result.ErrorOrNil()
}
