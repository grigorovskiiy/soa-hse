package postsservice

import (
	"context"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/models"
	pb "github.com/grigorovskiiy/soa-hse/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type PostsRepository interface {
	CreatePost(*models.DbPost) error
	DeletePost(int32, int32) error
	UpdatePost(*models.DbPost) error
	GetPost(int32, int32) (*models.DbPost, error)
	GetPostList(int32, int32, int32) ([]*models.DbPost, error)
	PostComment(*models.DbComment) error
	PostLike(*models.DbLike) error
	PostView(*models.DbView) error
	GetCommentList(int32, int32, int32) ([]*models.DbComment, error)
}
type Service struct {
	repository PostsRepository
}

func NewService(repository PostsRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreatePost(_ context.Context, pb *pb.PostDataRequest, userID int32) error {
	post := models.DbPost{
		Name:         pb.PostName,
		Tags:         pb.Tags,
		Description:  pb.PostDescription,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UserId:       int(userID),
		SecurityFlag: pb.SecurityFlag,
	}

	if err := s.repository.CreatePost(&post); err != nil {
		logger.Logger.Error("create post error", "error", err.Error())
		return err
	}

	return nil
}

func (s *Service) DeletePost(_ context.Context, pb *pb.PostID, userID int32) error {
	if err := s.repository.DeletePost(pb.PostId, userID); err != nil {
		logger.Logger.Error("delete post error", "error", err.Error())
		return err
	}

	return nil
}

func (s *Service) UpdatePost(_ context.Context, pb *pb.UpdatePostRequest, userID int32) error {
	post := models.DbPost{
		UpdatedAt:    time.Now(),
		Id:           int(pb.PostId),
		Name:         pb.PostData.PostName,
		Description:  pb.PostData.PostDescription,
		SecurityFlag: pb.PostData.SecurityFlag,
		Tags:         pb.PostData.Tags,
		UserId:       int(userID),
	}

	if err := s.repository.UpdatePost(&post); err != nil {
		logger.Logger.Error("update post error", "error", err.Error())
		return err
	}

	return nil
}

func (s *Service) GetPost(_ context.Context, p *pb.PostID, userID int32) (*pb.PostDataResponse, error) {
	postInfo, err := s.repository.GetPost(p.PostId, userID)
	if err != nil {
		logger.Logger.Error("get post info error", "error", err.Error())
		return nil, err
	}

	pb := pb.PostDataResponse{
		PostId:          int32(postInfo.Id),
		Tags:            postInfo.Tags,
		PostName:        postInfo.Name,
		PostDescription: postInfo.Description,
		SecurityFlag:    postInfo.SecurityFlag,
		CreatedAt:       timestamppb.New(postInfo.CreatedAt),
		UpdatedAt:       timestamppb.New(postInfo.UpdatedAt),
	}

	return &pb, nil
}

func (s *Service) GetPostList(_ context.Context, p *pb.PaginatedListRequest, userID int32) (*pb.ListPostsResponse, error) {
	posts, err := s.repository.GetPostList(p.Page, p.PageSize, userID)
	if err != nil {
		logger.Logger.Error("get post list error", "error", err.Error())
		return nil, err
	}
	pbPosts := pb.ListPostsResponse{
		Posts: make([]*pb.PostDataResponse, len(posts)),
	}
	for ind := range posts {
		pbPosts.Posts[ind] = &pb.PostDataResponse{
			PostId:          int32(posts[ind].Id),
			Tags:            posts[ind].Tags,
			PostName:        posts[ind].Name,
			PostDescription: posts[ind].Description,
			SecurityFlag:    posts[ind].SecurityFlag,
			CreatedAt:       timestamppb.New(posts[ind].CreatedAt),
			UpdatedAt:       timestamppb.New(posts[ind].UpdatedAt),
			UserId:          int32(posts[ind].UserId),
		}
	}

	return &pbPosts, nil
}

func (s *Service) PostComment(ctx context.Context, pb *pb.PostCommentRequest, userID int32) error {
	comment := models.DbComment{
		PostId:      int(pb.PostId),
		UserId:      int(userID),
		Description: pb.CommentDescription,
	}

	if err := s.repository.PostComment(&comment); err != nil {
		logger.Logger.Error("post comment error", "error", err.Error())
		return err
	}

	return nil
}

func (s *Service) PostLike(ctx context.Context, pb *pb.PostID, userID int32) error {
	like := models.DbLike{
		PostId: int(pb.PostId),
		UserId: int(userID),
	}

	if err := s.repository.PostLike(&like); err != nil {
		logger.Logger.Error("post like error", "error", err.Error())
		return err
	}

	return nil
}

func (s *Service) PostView(ctx context.Context, pb *pb.PostID, userID int32) error {
	view := models.DbView{
		PostId: int(pb.PostId),
		UserId: int(userID),
	}

	if err := s.repository.PostView(&view); err != nil {
		logger.Logger.Error("post view error", "error", err.Error())
		return err
	}

	return nil
}

func (s *Service) GetCommentList(_ context.Context, p *pb.PaginatedListRequest, userID int32) (*pb.ListCommentsResponse, error) {
	comments, err := s.repository.GetCommentList(p.Page, p.PageSize, userID)
	if err != nil {
		logger.Logger.Error("get comment list error", "error", err.Error())
		return nil, err
	}

	pbComments := pb.ListCommentsResponse{
		Comments: make([]*pb.CommentDataResponse, len(comments)),
	}
	for ind := range comments {
		pbComments.Comments[ind] = &pb.CommentDataResponse{
			CommentId:          int32(comments[ind].Id),
			PostId:             int32(comments[ind].PostId),
			CommentDescription: comments[ind].Description,
		}
	}

	return &pbComments, nil
}
