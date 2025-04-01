package server

import (
	"auth/posts_service/internal/application"
	"auth/posts_service/internal/infrastructure"
	pb "auth/protos"
	"context"
	"errors"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net"
	"os"
)

func NewServer(s *application.PostsServiceServer) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost%s", os.Getenv("POST_SERVICE_PORT")))
	if err != nil {
		infrastructure.Logger.Error(fmt.Sprintf("failed to listen: %s", err.Error()))
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
