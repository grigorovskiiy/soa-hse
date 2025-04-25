package main

import (
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/application"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/db"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/repository"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/server"
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
		fx.Provide(postsservice.NewPService),
		fx.Provide(func(s *postsservice.PService) application.PostsService {
			return s
		}),
		fx.Provide(application.NewPostsServer),
		fx.Provide(server.NewServer),
		fx.Invoke(server.RunServer),
	)

	fx.New(addOpts).Run()
}
