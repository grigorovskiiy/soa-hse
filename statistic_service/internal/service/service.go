package service

import (
	"context"
	pb "github.com/grigorovskiiy/soa-hse/protos"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StatisticRepository interface {
	GetViewsCount(ctx context.Context, postID int) (int, error)
	GetCommentsCount(ctx context.Context, postID int) (int, error)
	GetLikesCount(ctx context.Context, postID int) (int, error)
	GetViewsDynamic(ctx context.Context, postID int) ([]*models.Dynamic, error)
	GetCommentsDynamic(ctx context.Context, postID int) ([]*models.Dynamic, error)
	GetLikesDynamic(ctx context.Context, postID int) ([]*models.Dynamic, error)
	GetTopTenPosts(ctx context.Context, par string) ([]int, error)
	GetTopTenUsers(ctx context.Context, par string) ([]int, error)
}

type Transactor interface {
	WithTransaction(context.Context, func(context.Context) error) error
	WithTransactionWithValue(context.Context, func(context.Context) (any, error)) (any, error)
}
type Service struct {
	repository StatisticRepository
	tr         Transactor
}

func NewService(repository StatisticRepository, tr Transactor) *Service {
	return &Service{
		repository: repository,
		tr:         tr,
	}
}

func (s *Service) GetViewsCount(ctx context.Context, p *pb.PostID) (*pb.CountResponse, error) {
	count, err := s.tr.WithTransactionWithValue(ctx, func(ctx context.Context) (any, error) {
		count, err := s.repository.GetViewsCount(ctx, int(p.PostId))
		if err != nil {
			logger.Logger.Error("get views count error", "error", err.Error())
			return nil, err
		}

		return count, nil
	})

	if err != nil {
		logger.Logger.Error("get views count error", "error", err.Error())
		return nil, err
	}

	return &pb.CountResponse{Count: int32(count.(int))}, nil
}

func (s *Service) GetCommentsCount(ctx context.Context, p *pb.PostID) (*pb.CountResponse, error) {
	count, err := s.tr.WithTransactionWithValue(ctx, func(ctx context.Context) (any, error) {
		count, err := s.repository.GetCommentsCount(ctx, int(p.PostId))
		if err != nil {
			logger.Logger.Error("get comments count error", "error", err.Error())
			return nil, err
		}

		return count, nil
	})

	if err != nil {
		logger.Logger.Error("get comments count error", "error", err.Error())
		return nil, err
	}

	return &pb.CountResponse{Count: int32(count.(int))}, nil
}

func (s *Service) GetLikesCount(ctx context.Context, p *pb.PostID) (*pb.CountResponse, error) {
	count, err := s.tr.WithTransactionWithValue(ctx, func(ctx context.Context) (any, error) {
		count, err := s.repository.GetLikesCount(ctx, int(p.PostId))
		if err != nil {
			logger.Logger.Error("get likes count error", "error", err.Error())
			return nil, err
		}

		return count, nil
	})

	if err != nil {
		logger.Logger.Error("get likes count error", "error", err.Error())
		return nil, err
	}

	return &pb.CountResponse{Count: int32(count.(int))}, nil
}

func (s *Service) GetViewsDynamic(ctx context.Context, p *pb.PostID) (*pb.DynamicListResponse, error) {
	dbDyn, err := s.tr.WithTransactionWithValue(ctx, func(ctx context.Context) (any, error) {
		dbDyn, err := s.repository.GetViewsDynamic(ctx, int(p.PostId))
		if err != nil {
			logger.Logger.Error("get views dynamic error", "error", err.Error())
			return nil, err
		}

		return dbDyn, nil
	})
	if err != nil {
		logger.Logger.Error("get views dynamic error", "error", err.Error())
		return nil, err
	}

	dyn := dbDyn.([]*models.Dynamic)
	pbDyn := pb.DynamicListResponse{Dynamic: make([]*pb.DynamicResponse, len(dyn))}
	for i := range dyn {
		pbDyn.Dynamic[i] = &pb.DynamicResponse{
			Count: &pb.CountResponse{Count: int32(dyn[i].Count)}, Data: timestamppb.New(dyn[i].Date),
		}
	}

	return &pbDyn, nil
}

func (s *Service) GetCommentsDynamic(ctx context.Context, p *pb.PostID) (*pb.DynamicListResponse, error) {
	dbDyn, err := s.tr.WithTransactionWithValue(ctx, func(ctx context.Context) (any, error) {
		dbDyn, err := s.repository.GetCommentsDynamic(ctx, int(p.PostId))
		if err != nil {
			logger.Logger.Error("get comments dynamic error", "error", err.Error())
			return nil, err
		}

		return dbDyn, nil
	})
	if err != nil {
		logger.Logger.Error("get comments dynamic error", "error", err.Error())
		return nil, err
	}

	dyn := dbDyn.([]*models.Dynamic)
	pbDyn := pb.DynamicListResponse{Dynamic: make([]*pb.DynamicResponse, len(dyn))}
	for i := range dyn {
		pbDyn.Dynamic[i] = &pb.DynamicResponse{
			Count: &pb.CountResponse{Count: int32(dyn[i].Count)}, Data: timestamppb.New(dyn[i].Date),
		}
	}

	return &pbDyn, nil
}

func (s *Service) GetLikesDynamic(ctx context.Context, p *pb.PostID) (*pb.DynamicListResponse, error) {
	dbDyn, err := s.tr.WithTransactionWithValue(ctx, func(ctx context.Context) (any, error) {
		dbDyn, err := s.repository.GetLikesDynamic(ctx, int(p.PostId))
		if err != nil {
			logger.Logger.Error("get likes dynamic error", "error", err.Error())
			return nil, err
		}

		return dbDyn, nil
	})
	if err != nil {
		logger.Logger.Error("get likes dynamic error", "error", err.Error())
		return nil, err
	}

	dyn := dbDyn.([]*models.Dynamic)
	pbDyn := pb.DynamicListResponse{Dynamic: make([]*pb.DynamicResponse, len(dyn))}
	for i := range dyn {
		pbDyn.Dynamic[i] = &pb.DynamicResponse{
			Count: &pb.CountResponse{Count: int32(dyn[i].Count)}, Data: timestamppb.New(dyn[i].Date),
		}
	}

	return &pbDyn, nil
}

func (s *Service) GetTopTenPosts(ctx context.Context, p *pb.TopTenParameter) (*pb.TopTenPostsResponse, error) {
	dbPosts, err := s.tr.WithTransactionWithValue(ctx, func(ctx context.Context) (any, error) {
		dbPosts, err := s.repository.GetTopTenPosts(ctx, p.GetPar())
		if err != nil {
			logger.Logger.Error("get top ten posts error", "error", err.Error())
			return nil, err
		}

		return dbPosts, nil
	})
	if err != nil {
		logger.Logger.Error("get top ten posts error", "error", err.Error())
		return nil, err
	}

	posts := dbPosts.([]int)
	pbPosts := pb.TopTenPostsResponse{Posts: make([]*pb.PostID, len(posts))}
	for i := range posts {
		pbPosts.Posts[i] = &pb.PostID{PostId: int32(posts[i])}
	}

	return &pbPosts, nil
}

func (s *Service) GetTopTenUsers(ctx context.Context, p *pb.TopTenParameter) (*pb.TopTenUsersResponse, error) {
	dbUsers, err := s.tr.WithTransactionWithValue(ctx, func(ctx context.Context) (any, error) {
		dbUsers, err := s.repository.GetTopTenUsers(ctx, p.GetPar())
		if err != nil {
			logger.Logger.Error("get top ten users error", "error", err.Error())
			return nil, err
		}

		return dbUsers, nil
	})
	if err != nil {
		logger.Logger.Error("get top ten users error", "error", err.Error())
		return nil, err
	}

	users := dbUsers.([]int)
	pbUsers := pb.TopTenUsersResponse{Users: make([]*pb.UserID, len(users))}
	for i := range users {
		pbUsers.Users[i] = &pb.UserID{UserId: int32(users[i])}
	}

	return &pbUsers, nil
}
