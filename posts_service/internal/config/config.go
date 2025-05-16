package config

import (
	"github.com/caarlos0/env/v8"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/logger"
)

type Config struct {
	KafkaConfig
	PostsServiceConfig
}

type KafkaConfig struct {
	Brokers       []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"kafka:19092"`
	CommentsTopic string   `env:"KAFKA_COMMENTS_TOPIC" envDefault:"comments.topic"`
	LikesTopic    string   `env:"KAFKA_LIKES_TOPIC" envDefault:"likes.topic"`
	ViewsTopic    string   `env:"KAFKA_VIEWS_TOPIC" envDefault:"views.topic"`
}

type PostsServiceConfig struct {
	PostsServicePort      string `env:"POST_SERVICE_PORT" envDefault:":50051"`
	PostsServiceHost      string `env:"POST_SERVICE_HOST" envDefault:"posts-service"`
	PostsPostgresDb       string `env:"POSTS_POSTGRES_DB" envDefault:"posts_db"`
	PostsPostgresUser     string `env:"POSTS_POSTGRES_USER" envDefault:"username"`
	PostsPostgresPassword string `env:"POSTS_POSTGRES_PASSWORD" envDefault:"password"`
	PostsPostgresPort     string `env:"POSTS_POSTGRES_PORT" envDefault:":5432"`
	PostsPostgresHost     string `env:"POSTS_POSTGRES_HOST" envDefault:"posts-postgres"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		logger.Logger.Error("error parsing config", "error", err.Error())
		return nil, err
	}

	return &cfg, nil
}
