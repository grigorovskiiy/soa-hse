package config

import (
	"github.com/caarlos0/env/v8"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure/logger"
)

type Config struct {
	PostsServiceConfig
	StatisticServiceConfig
	UsersServiceConfig
	GatewayServiceConfig
}

type PostsServiceConfig struct {
	PostsServicePort string `env:"POST_SERVICE_PORT" envDefault:":50051"`
	PostsServiceHost string `env:"POST_SERVICE_HOST" envDefault:"posts-service"`
}

type StatisticServiceConfig struct {
	StatisticServicePort string `env:"STATISTIC_SERVICE_PORT" envDefault:":50052"`
	StatisticServiceHost string `env:"STATISTIC_SERVICE_HOST" envDefault:"statistic-service"`
}

type UsersServiceConfig struct {
	UsersServicePort string `env:"USERS_SERVICE_PORT" envDefault:":8081"`
	UsersServiceHost string `env:"USERS_SERVICE_HOST" envDefault:"users-service"`
}

type GatewayServiceConfig struct {
	GatewayServicePort string `env:"GATEWAY_SERVICE_PORT" envDefault:":8080"`
	GatewayServiceHost string `env:"GATEWAY_SERVICE_HOST" envDefault:"api-gateway-service"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		logger.Logger.Error("error parsing config", "error", err.Error())
		return nil, err
	}

	return &cfg, nil
}
