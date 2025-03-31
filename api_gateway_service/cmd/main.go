package main

import (
	_ "auth/api_gateway_service/docs"
	"auth/api_gateway_service/internal/api"
	"auth/api_gateway_service/internal/application"
	"auth/api_gateway_service/internal/infrastructure"
	"auth/api_gateway_service/internal/server"
	"github.com/joho/godotenv"

	"go.uber.org/fx"
)

func init() {
	if err := godotenv.Load(); err != nil {
		infrastructure.Logger.Error("env file is not found")
	}
}

// @title Swagger  API Gateway Service
// @version 1.0
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	addOpts := fx.Options(
		fx.Provide(application.NewGatewayApp),
		fx.Provide(api.NewApiGatewayService),
		fx.Provide(server.NewServer),
		fx.Invoke(server.RunServer),
		fx.Invoke(infrastructure.InitLogger),
	)

	fx.New(addOpts).Run()
}
