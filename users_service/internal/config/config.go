package config

import (
	"github.com/caarlos0/env/v8"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/logger"
)

type Config struct {
	KafkaConfig
	UsersServiceConfig
}

type KafkaConfig struct {
	Brokers      []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"kafka:19092"`
	ClientsTopic string   `env:"KAFKA_CLIENTS_TOPIC" envDefault:"clients.topic"`
}

type UsersServiceConfig struct {
	UsersServicePort      string `env:"USERS_SERVICE_PORT" envDefault:":8081"`
	UsersServiceHost      string `env:"USERS_SERVICE_HOST" envDefault:"users-service"`
	UsersPostgresDb       string `env:"USERS_POSTGRES_DB" envDefault:"users_db"`
	UsersPostgresUser     string `env:"USERS_POSTGRES_USER" envDefault:"username"`
	UsersPostgresPassword string `env:"USERS_POSTGRES_PASSWORD" envDefault:"password"`
	UsersPostgresPort     string `env:"USERS_POSTGRES_PORT" envDefault:":5432"`
	UsersPostgresHost     string `env:"USERS_POSTGRES_HOST" envDefault:"users-postgres"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		logger.Logger.Error("error parsing config", "error", err.Error())
		return nil, err
	}

	return &cfg, nil
}
