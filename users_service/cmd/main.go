package main

import (
	"auth/users_service/internal/api"
	"auth/users_service/internal/application"
	"auth/users_service/internal/infrastructure"
	"auth/users_service/internal/infrastructure/db"
	"auth/users_service/internal/infrastructure/repository"
	"auth/users_service/internal/server"
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
		fx.Provide(repository.NewUsersRepository),
		fx.Provide(api.NewUsersService),
		fx.Provide(application.NewUsersApp),
		fx.Provide(server.NewServer),
		fx.Invoke(server.RunServer),
	)

	fx.New(addOpts).Run()
}
