package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/application"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/logger"
	pb "github.com/grigorovskiiy/soa-hse/protos"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net"
)

func NewServer(s *application.PostsServiceApp, cfg *config.Config) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s%s", cfg.PostsServiceHost, cfg.PostsServicePort))
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("failed to listen: %s", err.Error()))
		return nil, nil
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPostsServiceServer(grpcServer, s)

	return grpcServer, lis
}

func RunServer(lc fx.Lifecycle, grpcServer *grpc.Server, listener net.Listener) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})

	return nil
}
