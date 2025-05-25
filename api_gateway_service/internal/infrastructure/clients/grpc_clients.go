package clients

import (
	"context"
	"fmt"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure/logger"
	pb "github.com/grigorovskiiy/soa-hse/protos"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	PostsServiceClient     pb.PostsServiceClient
	StatisticServiceClient pb.StatisticServiceClient
}

func NewGRPCClients(lc fx.Lifecycle, cfg *config.Config) (*GRPCClients, error) {
	postsConn, err := grpc.NewClient(fmt.Sprintf("%s%s", cfg.PostsServiceHost, cfg.PostsServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logger.Logger.Error("error creating posts service grpc client", err.Error())
		return nil, err
	}

	statisticConn, err := grpc.NewClient(fmt.Sprintf("%s%s", cfg.StatisticServiceHost, cfg.StatisticServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logger.Logger.Error("error creating statistic service grpc client", err.Error())
		return nil, err
	}

	clients := &GRPCClients{
		PostsServiceClient:     pb.NewPostsServiceClient(postsConn),
		StatisticServiceClient: pb.NewStatisticServiceClient(statisticConn),
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err := postsConn.Close(); err != nil {
				logger.Logger.Error("error closing posts service grpc client", err.Error())
				return err
			}
			if err := statisticConn.Close(); err != nil {
				logger.Logger.Error("error closing statistic service grpc client", err.Error())
				return err
			}

			return nil
		},
	})

	return clients, nil
}
