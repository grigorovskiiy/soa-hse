package application

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure"
	pb "github.com/grigorovskiiy/soa-hse/protos"
	"google.golang.org/grpc/metadata"
	"strconv"
)

type PostsService interface {
	CreatePost(*pb.PostDataRequest, int32) error
	DeletePost(*pb.PostID, int32) error
	UpdatePost(*pb.UpdatePostRequest, int32) error
	GetPost(*pb.PostID, int32) (*pb.PostDataResponse, error)
	GetPostList(*pb.ListPostsRequest, int32) (*pb.ListPostsResponse, error)
}

type PostsServiceServer struct {
	pb.UnimplementedPostsServiceServer
	Service PostsService
}

func NewPostsServer(service PostsService) *PostsServiceServer {
	return &PostsServiceServer{Service: service}
}

func (s *PostsServiceServer) CreatePost(ctx context.Context, pb *pb.PostDataRequest) (*empty.Empty, error) {
	logger := infrastructure.Logger.With("method", "CreatePost")
	logger.Info("posts grpc request started")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("no metadata from incoming context")
		return nil, errors.New("no metadata from incoming context")
	}
	m := md.Get("user_id")[0]
	userId, err := strconv.Atoi(m)
	if err != nil {
		logger.Error("wrong metadata", "error", err.Error())
		return nil, err
	}

	err = s.Service.CreatePost(pb, int32(userId))
	if err != nil {
		logger.Error("create post error", "error", err.Error())
		return nil, err
	}

	logger.Info("posts grpc request completed")

	return &empty.Empty{}, nil

}

func (s *PostsServiceServer) DeletePost(ctx context.Context, pb *pb.PostID) (*empty.Empty, error) {
	logger := infrastructure.Logger.With("method", "DeletePost")
	logger.Info("posts grpc request started")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("no metadata from incoming context")
		return nil, errors.New("no metadata from incoming context")
	}
	m := md.Get("user_id")[0]
	userId, err := strconv.Atoi(m)
	if err != nil {
		logger.Error("wrong metadata", "error", err.Error())
		return nil, err
	}

	err = s.Service.DeletePost(pb, int32(userId))
	if err != nil {
		logger.Error("delete post error", "error", err.Error())
		return nil, err
	}
	logger.Info("posts grpc request completed")
	return &empty.Empty{}, nil
}

func (s *PostsServiceServer) UpdatePost(ctx context.Context, pb *pb.UpdatePostRequest) (*empty.Empty, error) {
	logger := infrastructure.Logger.With("method", "UpdatePost")
	logger.Info("posts grpc request started")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("no metadata from incoming context")
		return nil, errors.New("no metadata from incoming context")
	}
	m := md.Get("user_id")[0]
	userId, err := strconv.Atoi(m)
	if err != nil {
		logger.Error("wrong metadata", "error", err.Error())
		return nil, err
	}

	err = s.Service.UpdatePost(pb, int32(userId))
	if err != nil {
		logger.Error("update post error", "error", err.Error())
		return nil, err
	}
	logger.Info("posts grpc request completed")

	return &empty.Empty{}, nil
}

func (s *PostsServiceServer) GetPost(ctx context.Context, pb *pb.PostID) (*pb.PostDataResponse, error) {
	logger := infrastructure.Logger.With("method", "GetPost")
	logger.Info("posts grpc request started")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("no metadata from incoming context")
		return nil, errors.New("no metadata from incoming context")
	}
	m := md.Get("user_id")[0]
	userId, err := strconv.Atoi(m)
	if err != nil {
		logger.Error("wrong metadata", "error", err.Error())
		return nil, err
	}

	postInfo, err := s.Service.GetPost(pb, int32(userId))
	if err != nil {
		logger.Error("get post request error", "error", err.Error())
		return nil, err
	}

	logger.Info("posts grpc request completed")
	return postInfo, nil
}

func (s *PostsServiceServer) GetPostList(ctx context.Context, pb *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	logger := infrastructure.Logger.With("method", "GetPostList")
	logger.Info("posts grpc request started")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("no metadata from incoming context")
		return nil, errors.New("no metadata from incoming context")
	}

	m := md.Get("user_id")[0]
	userId, err := strconv.Atoi(m)
	if err != nil {
		logger.Error("wrong metadata", "error", err.Error())
		return nil, err
	}
	postList, err := s.Service.GetPostList(pb, int32(userId))
	if err != nil {
		logger.Error("get post list error", "error", err.Error())
		return nil, err
	}
	logger.Info("posts grpc request completed")
	return postList, nil
}
