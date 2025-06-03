package main

import (
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/application"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/db"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/kafka"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/repository"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/server"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/service/eventsservice"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/service/postsservice"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Logger.Error("env file is not found")
	}
}

func main() {
	addOpts := fx.Options(
		fx.Provide(config.NewConfig),
		fx.Provide(db.InitDb),
		fx.Provide(repository.NewPRepository),
		fx.Provide(func(r *repository.PRepository) postsservice.PostsRepository {
			return r
		}),
		fx.Provide(postsservice.NewService),
		fx.Provide(func(s *postsservice.Service) application.PostsService {
			return s
		}),
		fx.Provide(kafka.NewBaseProducer),
		fx.Provide(eventsservice.NewKafkaService),
		fx.Provide(func(s *eventsservice.KafkaService) application.EventsService {
			return s
		}),
		fx.Provide(application.NewPostsApp),
		fx.Provide(server.NewServer),
		fx.Invoke(server.RunServer),
	)

	fx.New(addOpts).Run()
}
