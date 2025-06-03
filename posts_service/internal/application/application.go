package application

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/models"
	pb "github.com/grigorovskiiy/soa-hse/protos"
	"google.golang.org/grpc/metadata"
	"strconv"
	"time"
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

type EventsService interface {
	SendEvent(context.Context, string, any) error
}

type PostsServiceApp struct {
	pb.UnimplementedPostsServiceServer
	PostsService  PostsService
	EventsService EventsService
	cfg           *config.Config
}

func NewPostsApp(pS PostsService, eS EventsService, cfg *config.Config) *PostsServiceApp {
	return &PostsServiceApp{PostsService: pS, EventsService: eS, cfg: cfg}
}

func (s *PostsServiceApp) CreatePost(ctx context.Context, pb *pb.PostDataRequest) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "CreatePost")
	logger.Info("posts grpc request started")
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}

	if err = s.PostsService.CreatePost(ctx, pb, userID); err != nil {
		logger.Error("create post error", "error", err.Error())
		return nil, err
	}

	logger.Info("posts grpc request completed")

	return &empty.Empty{}, nil
}

func (s *PostsServiceApp) DeletePost(ctx context.Context, pb *pb.PostID) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "DeletePost")
	logger.Info("posts grpc request started")
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}

	if err = s.PostsService.DeletePost(ctx, pb, userID); err != nil {
		logger.Error("delete post error", "error", err.Error())
		return nil, err
	}

	logger.Info("posts grpc request completed")
	return &empty.Empty{}, nil
}

func (s *PostsServiceApp) UpdatePost(ctx context.Context, pb *pb.UpdatePostRequest) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "UpdatePost")
	logger.Info("posts grpc request started")
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}

	if err = s.PostsService.UpdatePost(ctx, pb, userID); err != nil {
		logger.Error("update post error", "error", err.Error())
		return nil, err
	}

	logger.Info("posts grpc request completed")
	return &empty.Empty{}, nil
}

func (s *PostsServiceApp) GetPost(ctx context.Context, pb *pb.PostID) (*pb.PostDataResponse, error) {
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

func (s *PostsServiceApp) GetPostList(ctx context.Context, pb *pb.PaginatedListRequest) (*pb.ListPostsResponse, error) {
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

func (s *PostsServiceApp) PostComment(ctx context.Context, pb *pb.PostCommentRequest) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "PostComment")
	logger.Info("posts grpc request started")

	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}

	if err = s.PostsService.PostComment(ctx, pb, userID); err != nil {
		logger.Error("post comment error", "error", err.Error())
		return nil, err
	}

	if err = s.EventsService.SendEvent(ctx, s.cfg.CommentsTopic, models.LikeViewCommentUpdate{PostId: int(pb.PostId), UserId: int(userID), Time: time.Now()}); err != nil {
		logger.Error("post comment send event error", "error", err.Error())
		return nil, err
	}

	logger.Info("posts grpc request completed")
	return &empty.Empty{}, nil
}

func (s *PostsServiceApp) PostLike(ctx context.Context, pb *pb.PostID) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "PostLike")
	logger.Info("posts grpc request started")
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}

	if err = s.PostsService.PostLike(ctx, pb, userID); err != nil {
		logger.Error("post like error", "error", err.Error())
		return nil, err
	}

	if err = s.EventsService.SendEvent(ctx, s.cfg.LikesTopic, models.LikeViewCommentUpdate{PostId: int(pb.PostId), UserId: int(userID), Time: time.Now()}); err != nil {
		logger.Error("post like send event error", "error", err.Error())
		return nil, err

	}

	logger.Info("posts grpc request completed")
	return &empty.Empty{}, nil
}

func (s *PostsServiceApp) PostView(ctx context.Context, pb *pb.PostID) (*empty.Empty, error) {
	logger := logger.Logger.With("method", "PostView")
	logger.Info("posts grpc request started")
	userID, err := GetUserID(ctx)
	if err != nil {
		logger.Error("error getting userID from ctx", "error", err.Error())
		return nil, err
	}

	if err = s.PostsService.PostView(ctx, pb, userID); err != nil {
		logger.Error("post view error", "error", err.Error())
		return nil, err
	}

	if err = s.EventsService.SendEvent(ctx, s.cfg.ViewsTopic, models.LikeViewCommentUpdate{PostId: int(pb.PostId), UserId: int(userID), Time: time.Now()}); err != nil {
		logger.Error("post view send event error", "error", err.Error())
		return nil, err
	}

	logger.Info("posts grpc request completed")
	return &empty.Empty{}, nil
}

func (s *PostsServiceApp) GetCommentList(ctx context.Context, pb *pb.PaginatedListRequest) (*pb.ListCommentsResponse, error) {
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
