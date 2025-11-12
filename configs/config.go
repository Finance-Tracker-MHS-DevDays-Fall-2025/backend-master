package configs

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServiceConfig struct {
	ServerCfg   ServerConfig
	DatabaseCfg DatabaseConfig
	SlavesCfg   SlavesConfig
}

type ServerConfig struct {
	GrpcPort string `env:"GRPC_PORT" env-required:"true"`
	HttpPort string `env:"HTTP_PORT" env-required:"true"`
}

type DatabaseConfig struct {
	PgHost string `env:"PG_HOST" env-required:"true"`
	PgPort string `env:"PG_PORT" env-required:"true"`
	PgDb   string `env:"PG_DB" env-required:"true"`
	PgUser string `env:"PG_USER" env-required:"true"`
	PgPass string `env:"PG_PASS" env-required:"true"`
}

type SlavesConfig struct {
	AnalyzerUrl     string `env:"ANALYZER_URL" env-required:"true"`
	MarketUrl       string `env:"MARKET_URL" env-required:"true"`
	WalletUrl       string `env:"WALLET_URL" env-required:"true"`
	NotificationUrl string `env:"NOTIFICATION_URL" env-required:"true"`
}

func New() (*ServiceConfig, error) {
	var cfg ServiceConfig

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("error while reading env: %w", err)
	}

	return &cfg, nil
}
