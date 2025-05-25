package main

import (
	_ "github.com/grigorovskiiy/soa-hse/api_gateway_service/docs"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/application"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure/clients"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/server"
	"github.com/joho/godotenv"

	"go.uber.org/fx"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Logger.Error("env file is not found")
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
		fx.Provide(
			config.NewConfig,
			clients.NewGRPCClients,
			application.NewGatewayApp,
			server.NewServer,
		),
		fx.Invoke(
			server.RunServer,
		),
	)

	fx.New(addOpts).Run()
}
