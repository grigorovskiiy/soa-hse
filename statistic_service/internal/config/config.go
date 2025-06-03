package config

import (
	"github.com/caarlos0/env/v8"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/logger"
)

type Config struct {
	ClickHouseConfig
	StatisticServiceServerConfig
}

type ClickHouseConfig struct {
	ClickHousePort     string `env:"CLICKHOUSE_PORT" envDefault:":9000"`
	ClickHouseUser     string `env:"CLICKHOUSE_USER" envDefault:"username"`
	ClickHousePassword string `env:"CLICKHOUSE_PASSWORD" envDefault:"password"`
	ClickHouseDb       string `env:"CLICKHOUSE_DB" envDefault:"clickhouse_db"`
	ClickHouseHost     string `env:"CLICKHOUSE_HOST" envDefault:"clickhouse"`
}

type StatisticServiceServerConfig struct {
	StatisticServicePort string `env:"STATISTIC_SERVICE_PORT" envDefault:":50052"`
	StatisticServiceHost string `env:"STATISTIC_SERVICE_HOST" envDefault:"statistic-service"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		logger.Logger.Error("error parsing config", "error", err.Error())
		return nil, err
	}

	return &cfg, nil
}
