package main

import (
	"auth/posts_service/internal/application"
	"auth/posts_service/internal/infrastructure"
	"auth/posts_service/internal/infrastructure/db"
	"auth/posts_service/internal/infrastructure/repository"
	"auth/posts_service/internal/server"
	"auth/posts_service/internal/service"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func init() {
	if err := godotenv.Load(); err != nil {
		infrastructure.Logger.Error("env file is not found")
	}
}

func main() {
	addOpts := fx.Options(
		fx.Invoke(infrastructure.InitLogger),
		fx.Provide(db.InitDb),
		fx.Provide(repository.NewPRepository),
		fx.Provide(func(r *repository.PRepository) service.PostsRepository {
			return r
		}),
		fx.Provide(service.NewUService),
		fx.Provide(func(s *service.UService) application.PostsService {
			return s
		}),
		fx.Provide(application.NewPostsServer),
		fx.Provide(server.NewServer),
		fx.Invoke(server.RunServer),
	)

	fx.New(addOpts).Run()
}
