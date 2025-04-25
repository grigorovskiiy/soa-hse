package application

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/logger"
	pb "github.com/grigorovskiiy/soa-hse/protos"
	"google.golang.org/grpc/metadata"
	"strconv"
)

type PostsService interface {
	CreatePost(context.Context, *pb.PostDataRequest, int32) error
	DeletePost(context.Context, *pb.PostID, int32) error
	UpdatePost(context.Context, *pb.UpdatePostRequest, int32) error
	GetPost(context.Context, *pb.PostID, int32) (*pb.PostDataResponse, error)
	GetPostList(context.Context, *pb.PaginatedListRequest, int32) (*pb.ListPostsResponse, error)
	PostComment(context.Context, *pb.PostCommentRequest, int32) error
	PostLike(context.Context, *pb.PostID, int32) error
	PostView(context.Context, *pb.PostID, int32) error
	GetCommentList(context.Context, *pb.PaginatedListRequest, int32) (*pb.ListCommentsResponse, error)
}

type KafkaService interface {
	SendUpdatePost(context.Context, string, any) error
}

type PostsServiceServer struct {
	pb.UnimplementedPostsServiceServer
	PostsService PostsService
	KafkaService KafkaService
}

func NewPostsServer(pS PostsService, kS KafkaService) *PostsServiceServer {
	return &PostsServiceServer{PostsService: pS, KafkaService: kS}
}

func (s *PostsServiceServer) CreatePost(ctx context.Context, pb *pb.PostDataRequest) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "CreatePost")
	logger.Info("posts grpc request started")
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}

	err = s.PostsService.CreatePost(ctx, pb, userID)
	if err != nil {
		logger.Error("create post error", "error", err.Error())
		return nil, err
	}

	logger.Info("posts grpc request completed")

	return &empty.Empty{}, nil
}

func (s *PostsServiceServer) DeletePost(ctx context.Context, pb *pb.PostID) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "DeletePost")
	logger.Info("posts grpc request started")
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}

	err = s.PostsService.DeletePost(ctx, pb, userID)
	if err != nil {
		logger.Error("delete post error", "error", err.Error())
		return nil, err
	}
	logger.Info("posts grpc request completed")
	return &empty.Empty{}, nil
}

func (s *PostsServiceServer) UpdatePost(ctx context.Context, pb *pb.UpdatePostRequest) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "UpdatePost")
	logger.Info("posts grpc request started")
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}

	err = s.PostsService.UpdatePost(ctx, pb, userID)
	if err != nil {
		logger.Error("update post error", "error", err.Error())
		return nil, err
	}
	logger.Info("posts grpc request completed")

	return &empty.Empty{}, nil
}

func (s *PostsServiceServer) GetPost(ctx context.Context, pb *pb.PostID) (*pb.PostDataResponse, error) {
	logger := logger.Logger.With("method", "GetPost")
	logger.Info("posts grpc request started")

	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}

	postInfo, err := s.PostsService.GetPost(ctx, pb, userID)
	if err != nil {
		logger.Error("get post request error", "error", err.Error())
		return nil, err
	}

	logger.Info("posts grpc request completed")
	return postInfo, nil
}

func (s *PostsServiceServer) GetPostList(ctx context.Context, pb *pb.PaginatedListRequest) (*pb.ListPostsResponse, error) {
	logger := logger.Logger.With("method", "GetPostList")
	logger.Info("posts grpc request started")

	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}
	postList, err := s.PostsService.GetPostList(ctx, pb, userID)
	if err != nil {
		logger.Error("get post list error", "error", err.Error())
		return nil, err
	}
	logger.Info("posts grpc request completed")
	return postList, nil
}

func (s *PostsServiceServer) PostComment(ctx context.Context, pb *pb.PostCommentRequest) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "PostComment")
	logger.Info("posts grpc request started")

	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}

	err = s.PostsService.PostComment(ctx, pb, userID)
	if err != nil {
		logger.Error("post comment error", "error", err.Error())
		return nil, err
	}
	logger.Info("posts grpc request completed")
	return &empty.Empty{}, nil
}

func (s *PostsServiceServer) PostLike(ctx context.Context, pb *pb.PostID) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "PostLike")
	logger.Info("posts grpc request started")
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}
	err = s.PostsService.PostLike(ctx, pb, userID)
	if err != nil {
		logger.Error("post like error", "error", err.Error())
		return nil, err
	}
	logger.Info("posts grpc request completed")
	return &empty.Empty{}, nil
}

func (s *PostsServiceServer) PostView(ctx context.Context, pb *pb.PostID) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "PostView")
	logger.Info("posts grpc request started")
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}
	err = s.PostsService.PostView(ctx, pb, userID)
	if err != nil {
		logger.Error("post view error", "error", err.Error())
		return nil, err
	}
	logger.Info("posts grpc request completed")
	return &empty.Empty{}, nil
}

func (s *PostsServiceServer) GetCommentsList(ctx context.Context, pb *pb.PaginatedListRequest) (*pb.ListCommentsResponse, error) {
	logger := logger.Logger.With("method", "GetCommentsList")
	logger.Info("posts grpc request started")
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}
	comments, err := s.PostsService.GetCommentList(ctx, pb, userID)
	if err != nil {
		logger.Error("get comments list error", "error", err.Error())
		return nil, err
	}
	logger.Info("posts grpc request completed")
	return comments, nil
}

func GetUserID(ctx context.Context) (int32, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("no metadata from incoming context")
	}
	m := md.Get("user_id")[0]
	userId, err := strconv.Atoi(m)
	if err != nil {
		return 0, err
	}

	return int32(userId), nil
}
