package application

import (
	"context"
	pb "github.com/grigorovskiiy/soa-hse/protos"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/logger"
)

type StatisticService interface {
	GetViewsCount(ctx context.Context, pb *pb.PostID) (*pb.CountResponse, error)
	GetCommentsCount(ctx context.Context, pb *pb.PostID) (*pb.CountResponse, error)
	GetLikesCount(ctx context.Context, pb *pb.PostID) (*pb.CountResponse, error)
	GetViewsDynamic(ctx context.Context, pb *pb.PostID) (*pb.DynamicListResponse, error)
	GetCommentsDynamic(ctx context.Context, pb *pb.PostID) (*pb.DynamicListResponse, error)
	GetLikesDynamic(ctx context.Context, pb *pb.PostID) (*pb.DynamicListResponse, error)
	GetTopTenPosts(ctx context.Context, in *pb.TopTenParameter) (*pb.TopTenPostsResponse, error)
	GetTopTenUsers(ctx context.Context, in *pb.TopTenParameter) (*pb.TopTenUsersResponse, error)
}

type StatisticServiceApp struct {
	pb.UnimplementedStatisticServiceServer
	StatisticService StatisticService
}

func NewStatisticServiceApp(StatisticService StatisticService) *StatisticServiceApp {
	return &StatisticServiceApp{StatisticService: StatisticService}
}

func (s *StatisticServiceApp) GetViewsCount(ctx context.Context, pb *pb.PostID) (*pb.CountResponse, error) {
	logger := logger.Logger.With("method", "GetViewsCount")
	logger.Info("statistic grpc request started")

	count, err := s.StatisticService.GetViewsCount(ctx, pb)
	if err != nil {
		logger.Error("error getting views count", "error", err.Error())
		return nil, err
	}

	logger.Info("statistic grpc request completed")

	return count, nil
}

func (s *StatisticServiceApp) GetCommentsCount(ctx context.Context, pb *pb.PostID) (*pb.CountResponse, error) {
	logger := logger.Logger.With("method", "GetCommentsCount")
	logger.Info("statistic grpc request started")

	count, err := s.StatisticService.GetCommentsCount(ctx, pb)
	if err != nil {
		logger.Error("error getting comments count", "error", err.Error())
		return nil, err
	}

	logger.Info("statistic grpc request completed")

	return count, nil
}

func (s *StatisticServiceApp) GetLikesCount(ctx context.Context, pb *pb.PostID) (*pb.CountResponse, error) {
	logger := logger.Logger.With("method", "GetLikesCount")
	logger.Info("statistic grpc request started")

	count, err := s.StatisticService.GetLikesCount(ctx, pb)
	if err != nil {
		logger.Error("error getting likes count", "error", err.Error())
		return nil, err
	}

	logger.Info("statistic grpc request completed")

	return count, nil
}

func (s *StatisticServiceApp) GetViewsDynamic(ctx context.Context, pb *pb.PostID) (*pb.DynamicListResponse, error) {
	logger := logger.Logger.With("method", "GetViewsDynamic")
	logger.Info("statistic grpc request started")

	dynamic, err := s.StatisticService.GetViewsDynamic(ctx, pb)
	if err != nil {
		logger.Error("error getting views dynamic", "error", err.Error())
		return nil, err
	}

	logger.Info("statistic grpc request completed")

	return dynamic, nil
}

func (s *StatisticServiceApp) GetCommentsDynamic(ctx context.Context, pb *pb.PostID) (*pb.DynamicListResponse, error) {
	logger := logger.Logger.With("method", "GetCommentsDynamic")
	logger.Info("statistic grpc request started")

	dynamic, err := s.StatisticService.GetViewsDynamic(ctx, pb)
	if err != nil {
		logger.Error("error getting comments dynamic", "error", err.Error())
		return nil, err
	}

	logger.Info("statistic grpc request completed")

	return dynamic, nil
}

func (s *StatisticServiceApp) GetLikesDynamic(ctx context.Context, pb *pb.PostID) (*pb.DynamicListResponse, error) {
	logger := logger.Logger.With("method", "GetLikesDynamic")
	logger.Info("statistic grpc request started")

	dynamic, err := s.StatisticService.GetLikesDynamic(ctx, pb)
	if err != nil {
		logger.Error("error getting likes dynamic", "error", err.Error())
		return nil, err
	}

	logger.Info("statistic grpc request completed")

	return dynamic, nil
}

func (s *StatisticServiceApp) GetTopTenPosts(ctx context.Context, pb *pb.TopTenParameter) (*pb.TopTenPostsResponse, error) {
	logger := logger.Logger.With("method", "GetTopTenPosts")
	logger.Info("statistic grpc request started")

	posts, err := s.StatisticService.GetTopTenPosts(ctx, pb)
	if err != nil {
		logger.Error("error getting top ten posts", "error", err.Error())
		return nil, err
	}

	logger.Info("statistic grpc request completed")

	return posts, nil
}

func (s *StatisticServiceApp) GetTopTenUsers(ctx context.Context, pb *pb.TopTenParameter) (*pb.TopTenUsersResponse, error) {
	logger := logger.Logger.With("method", "GetTopTenPosts")
	logger.Info("statistic grpc request started")

	users, err := s.StatisticService.GetTopTenUsers(ctx, pb)
	if err != nil {
		logger.Error("error getting top ten posts", "error", err.Error())
		return nil, err
	}

	logger.Info("statistic grpc request completed")

	return users, nil
}
