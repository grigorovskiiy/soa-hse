package clients

import (
	"context"
	"fmt"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure"
	pb "github.com/grigorovskiiy/soa-hse/protos"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

type GRPCClients struct {
	PostsServiceClient pb.PostsServiceClient
}

func NewGRPCClients(lc fx.Lifecycle) (*GRPCClients, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("posts-service%s", os.Getenv("POST_SERVICE_PORT")),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		infrastructure.Logger.Error("error creating grpc clients connection", err.Error())
		return nil, err
	}

	clients := &GRPCClients{
		PostsServiceClient: pb.NewPostsServiceClient(conn),
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return clients, nil
}
