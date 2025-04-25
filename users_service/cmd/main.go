package main

import (
	"github.com/grigorovskiiy/soa-hse/users_service/internal/application"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/db"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/repository"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/server"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/service"
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
		fx.Provide(db.InitDb),
		fx.Provide(repository.NewUsersRepository),
		fx.Provide(func(r *repository.UsersRepository) service.Repository {
			return r
		}),
		fx.Provide(service.NewUService),
		fx.Provide(func(s *service.UService) application.UsersService {
			return s
		}),
		fx.Provide(application.NewUsersApp),
		fx.Provide(server.NewServer),
		fx.Invoke(server.RunServer),
	)

	fx.New(addOpts).Run()
}
