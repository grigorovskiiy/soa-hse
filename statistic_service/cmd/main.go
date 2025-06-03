package main

import (
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/application"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/db"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/repository"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/repository/txs"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/server"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/service"
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
		fx.Provide(
			config.NewConfig,
			db.InitDb,
			repository.NewRepository,
			func(r *repository.Repository) service.StatisticRepository {
				return r
			},
			txs.NewTxBeginner,
			func(t *txs.TxBeginner) service.Transactor {
				return t
			},
			service.NewService,
			func(s *service.Service) application.StatisticService {
				return s
			},
			application.NewStatisticServiceApp,
			server.NewServer,
		),
		fx.Invoke(server.RunServer),
	)

	fx.New(addOpts).Run()
}
